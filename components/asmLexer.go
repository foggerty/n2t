/*x
  asmLexer

  Input: a string

  Output: a channel of asmLexemes
*/

package components

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

////////////////////////////////////////////////////////////////////////////////
// Public

// StateFunction represents the lexer's current 'state'.  i.e. are we
// at the beginning of an A-instruction, a label or a C-Instruction.
type stateFunction func(*Lexer) stateFunction

// Lexer tracks the progress as the lexer process moves through the
// input string.
type Lexer struct {
	input    string         // entire source file, not bothering with streaming (for now)
	start    int            // start of current item in bytes, NOT characters
	pos      int            // the position as we search along/end of current item
	width    int            // width of last rune that was read
	items    chan AsmLexeme // channel on which to pass back the tokens
	lineNum  int            // current source line number
	haveComp bool           // true if we've just processed an instruction
}

// NewLexer returns both a lexer structure, and its output channel, on
// which 'lexenes' (is that an actual word?) get passed as they are
// read.
func newLexer(input string) chan AsmLexeme {
	l := &Lexer{
		input:   input,
		items:   make(chan AsmLexeme),
		lineNum: 1}

	go l.run()

	return l.items
}

// Run starts the lexer process.
func (l *Lexer) run() {
	defer close(l.items)

	state := initState(l)

	for state != nil {
		state = state(l)
	}
}

////////////////////////////////////////////////////////////////////////////////
// character sets for various tokens (symbol, instruction etc)

const validSymbol string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.$:-"
const validInstruction string = "-+01!AMD&|nullJGELMNTQEP"

////////////////////////////////////////////////////////////////////////////////
// helper routines

func (l *Lexer) accept(valid string) {
	next := l.nextRune()

	for strings.ContainsRune(valid, next) {
		next = l.nextRune()
	}

	l.rewind()
}

func (l *Lexer) skip(valid string) {
	l.start = l.pos
	next := l.peek()

	// it turns out that Go, C# and Ruby all think that ANY string
	// contains the empty string.  So I figure this must be pretty much
	// language-standard.  Doh :-)

	for strings.Contains(valid, next) && next != "" {
		l.skipOne()
		next = l.peek()
	}
}

// emit will thrown the value of pos-start from input onto the output
// channel
func (l *Lexer) emit(aI AsmInstruction) {
	if aI == asmEOL && !l.haveComp {
		l.lineNum++
		return
	}

	l.haveComp = aI != asmEOL

	var value string

	if aI == asmEOF || aI == asmEOL {
		value = ""
	} else {
		value = l.input[l.start:l.pos]
	}

	l.items <- AsmLexeme{
		instruction: aI,
		value:       value,
		lineNum:     l.lineNum}

	if aI == asmEOL {
		l.lineNum++
	}

	l.start = l.pos
	l.width = 0
}

func (l *Lexer) skipOne() {
	l.nextRune()
	l.start = l.pos
}

func (l *Lexer) skipEol() {
	n := l.peek()

	if n == "\r" {
		l.skipOne()
		l.skipOne()
	}

	if n == "\n" {
		l.skipOne()
	}
}

// peek peeks at the next rune, without advancing through the string.
func (l *Lexer) peek() string {
	if l.atEof() {
		return ""
	}

	result := l.nextRune()
	l.rewind()

	return string(result)
}

// peeks at the next x runes
func (l *Lexer) peekAt(x int) string {
	current := l.pos
	var result string

	for x > 0 {
		rune, width := utf8.DecodeRuneInString(l.input[current:])
		result += string(rune)
		current += width
		x--
	}

	return result
}

// nextRune returns the next rune in input, moving forward.
func (l *Lexer) nextRune() rune {
	rune, width := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += width
	l.width = width

	return rune
}

// Just emits an error lexeme.
func (l *Lexer) error() stateFunction {
	l.items <- AsmLexeme{
		instruction: asmERROR,
		value:       fmt.Sprintf("Unknown error at line %d (%s)", l.lineNum, l.input[l.start:l.pos])}

	return nil
}

// atEof ....  Pretty sure you can figure this one out.
func (l *Lexer) atEof() bool {
	return l.pos >= len(l.input)
}

// True if next rune is \n, otherwise false.
func (l *Lexer) atEol() bool {
	if l.atEof() {
		return false
	}

	next := l.peek()

	return next == "\r" || next == "\n"
}

func (l *Lexer) moveToEol() {
	// first check that there is an EOL
	i := strings.Index(l.input[l.pos:], "\r")

	if i == -1 {
		i = strings.Index(l.input[l.pos:], "\n")
	}

	// is this the last line?
	if i == -1 {
		// move to EOF
		i = len(l.input)
	} else {
		i += l.pos
	}

	l.start, l.pos = i, i
}

// rewind moves back to the previous rune.
func (l *Lexer) rewind() {
	l.pos -= l.width
	l.width = 0
}

// RewindTo rewinds BOTH pos and start to a set location.
func (l *Lexer) rewindTo(p int) error {
	if p < l.start {
		return errors.New("Attempted to rewind past start")
	}

	l.pos = p

	return nil
}

func (l *Lexer) atBeginComment() bool {
	if l.atEof() {
		return false
	}

	test := l.peekAt(2)

	return test == "//"
}

func (l *Lexer) noCapture() bool {
	return (l.pos - l.start) == 0
}

////////////////////////////////////////////////////////////////////////////////
// State functions

// General "move forward until we find something useful" function.
// Makes no assumptions about where it is.
func initState(l *Lexer) stateFunction {
	_ = "breakpoint"
	l.skip(" \t")

	if l.atEol() {
		l.skipEol()
		l.emit(asmEOL)

		return initState
	}

	if l.atEof() {
		l.emit(asmEOF)
		return nil
	}

	for l.atBeginComment() {
		l.moveToEol()
		return initState
	}

	next := l.peek()

	if next == "@" {
		l.skipOne()
		return aInstruction
	}

	if next == "(" {
		l.skipOne()
		return atLabel
	}

	return atCInstruction
}

// Moves to the EOL.
// Assumes that at the end of an instruction or a symbol, so
// anything other than white space or a comment before EOL or EOF
// is an error.
func endOfCode(l *Lexer) stateFunction {
	l.skip(" \t")

	if l.atEof() {
		l.emit(asmEOF)
		return nil
	}

	if l.atBeginComment() {
		l.moveToEol()
		return initState
	}

	next := l.peek()

	if next == "\r" || next == "\n" {
		return initState
	}

	return errState
}

// Handles A-instructions (@123 / @loop)
// Assumes that it is positioned just after the '@'.
func aInstruction(l *Lexer) stateFunction {
	l.accept(validSymbol)

	if l.noCapture() {
		return errState
	}

	l.emit(asmAINSTRUCT)

	return endOfCode
}

// Handles labels i.e. (LOOP)
// Assumes that it is positioned just after the opening '('
func atLabel(l *Lexer) stateFunction {
	l.accept(validSymbol)

	next := l.peek()

	if l.noCapture() || next != ")" {
		return errState
	}

	l.emit(asmLABEL)
	l.skipOne() // step over the closing ')'

	return endOfCode
}

// Handles C-Instructions.
// Note that there are four ways this can be represented:
//
//    d=c;j
//    d=c
//    c;j
//    c
//
// this will figure out in which situation we're in, and set the next
// state accordingly.
//
// Assumes that it's at the beginning of some text.  Hopefully
// it's a C-Instruction.
func atCInstruction(l *Lexer) stateFunction {
	l.accept(validInstruction)

	if l.noCapture() {
		return errState
	}

	next := l.peek()

	// destination part of d=c;j
	if next == "=" {
		l.emit(asmDEST)
		l.skipOne()
		return atCmp
	}

	// comp part of c;j
	if next == ";" {
		l.emit(asmCOMP)
		l.skipOne()
		return atJmp
	}

	// This may be a comp instruction, but by itself that does nothing
	// other than burn cycles.
	return errState
}

// Handles the 'c' part of a C-Instruction (d=c;j)
// Assumes positioned immediately after the '='.
func atCmp(l *Lexer) stateFunction {
	l.accept(validInstruction)

	if l.noCapture() {
		return errState
	}

	l.emit(asmCOMP)

	next := l.peek()

	if next == ";" {
		l.skipOne()
		return atJmp
	}

	return endOfCode
}

// Handles the 'j' part of a C-Instruction.
// Assumes is positioned immediately after the ';'
func atJmp(l *Lexer) stateFunction {
	l.accept(validInstruction)

	if l.noCapture() {
		return errState
	}

	l.emit(asmJUMP)

	return endOfCode
}

// Assumes that it's all gone to hell.
// Emit an error, and moves past the EOL, return skipWhitespace
// If there is no EOL before EOF, return nil.
func errState(l *Lexer) stateFunction {
	l.error()

	l.moveToEol()

	if l.atEof() {
		l.emit(asmEOF)
		return nil
	}

	return initState
}
