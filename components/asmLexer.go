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
type stateFunction func(*lexer) stateFunction

// Lexer tracks the progress as the lexer process moves through the
// input string.
type lexer struct {
	input    string         // entire source file, not bothering with streaming (for now)
	start    int            // start of current item in bytes, NOT characters
	pos      int            // the position as we search along/end of current item
	width    int            // width of last rune that was read
	items    chan asmLexeme // channel on which to pass back the tokens
	lineNum  int            // current source line number
	instruct bool           // true if we've started processing an instruction
}

// NewLexer returns both a lexer structure, and its output channel, on
// which 'lexenes' (is that an actual word?) get passed as they are
// read.
func newLexer(input string) chan asmLexeme {
	l := &lexer{
		input:   input,
		items:   make(chan asmLexeme),
		lineNum: 1}

	go l.run()

	return l.items
}

// Run starts the lexer process.
func (l *lexer) run() {
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

func (l *lexer) accept(valid string) {
	next := l.nextRune()

	for strings.ContainsRune(valid, next) {
		next = l.nextRune()
	}

	l.rewind()
}

func (l *lexer) skip(valid string) {
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
func (l *lexer) emit(aI asmInstruction) {
	// Want to collapse multiple EOLs, i.e. skip over empty lines and
	// comments.  So if we're not processing an instruction when we
	// receive an EOL, keep going but still increment the line count.

	if !l.instruct && aI == asmEOL {
		l.lineNum++
		return
	}

	var value string

	if aI == asmEOF || aI == asmEOL {
		value = ""
	} else {
		value = l.input[l.start:l.pos]
	}

	l.items <- asmLexeme{
		instruction: aI,
		value:       value,
		lineNum:     l.lineNum}

	if aI == asmEOL {
		l.lineNum++
	}

	l.start = l.pos
	l.width = 0
	l.instruct = aI != asmEOL
}

func (l *lexer) skipOne() {
	l.nextRune()
	l.start = l.pos
}

func (l *lexer) skipEol() {
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
func (l *lexer) peek() string {
	if l.atEof() {
		return ""
	}

	result := l.nextRune()
	l.rewind()

	return string(result)
}

// peeks at the next x runes
func (l *lexer) peekAt(x int) string {
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
func (l *lexer) nextRune() rune {
	rune, width := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += width
	l.width = width

	return rune
}

// Just emits an error lexeme.

func (l *lexer) error() {
	// just dump the entire line
	start := l.start

	for start > 0 {
		if string(l.input[start]) == "\n" {
			break
		}
		start--
	}

	end := strings.IndexAny(l.input[l.start:], "\r\n")

	if end == -1 {
		end = len(l.input)
	} else {
		end += l.start
	}

	badLine := strings.Trim(l.input[start:end], "\r\n\t ")

	l.items <- asmLexeme{
		instruction: asmERROR,
		value:       fmt.Sprintf("Unknown error at line %d: (%s)", l.lineNum, badLine),
		lineNum:     l.lineNum}
}

// atEof ....  Pretty sure you can figure this one out.
func (l *lexer) atEof() bool {
	return l.pos >= len(l.input)
}

// True if next rune is \n, otherwise false.
func (l *lexer) atEol() bool {
	if l.atEof() {
		return false
	}

	next := l.peek()

	return next == "\r" || next == "\n"
}

func (l *lexer) moveToEol() {
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
func (l *lexer) rewind() {
	l.pos -= l.width
	l.width = 0
}

// RewindTo rewinds BOTH pos and start to a set location.
func (l *lexer) rewindTo(p int) error {
	if p < l.start {
		return errors.New("Attempted to rewind past start")
	}

	l.pos = p

	return nil
}

func (l *lexer) atBeginComment() bool {
	if l.atEof() {
		return false
	}

	test := l.peekAt(2)

	return test == "//"
}

func (l *lexer) noCapture() bool {
	return (l.pos - l.start) == 0
}

////////////////////////////////////////////////////////////////////////////////
// State functions

// General "move forward until we find something useful" function.
// Makes no assumptions about where it is.
func initState(l *lexer) stateFunction {
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
func endOfCode(l *lexer) stateFunction {
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
func aInstruction(l *lexer) stateFunction {
	l.accept(validSymbol)

	if l.noCapture() {
		return errState
	}

	l.emit(asmAINSTRUCT)

	return endOfCode
}

// Handles labels i.e. (LOOP)
// Assumes that it is positioned just after the opening '('
func atLabel(l *lexer) stateFunction {
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
func atCInstruction(l *lexer) stateFunction {
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
func atCmp(l *lexer) stateFunction {
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
func atJmp(l *lexer) stateFunction {
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
func errState(l *lexer) stateFunction {
	l.error()

	l.moveToEol()

	if l.atEof() {
		l.emit(asmEOF)
		return nil
	}

	return initState
}
