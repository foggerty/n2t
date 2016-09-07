/*

Package components has all of the components requires for both the
assembler and compiler that are written for the Nand2Tetris course
(http://nand2tetris.org/)

Idea blatantly stolen from Rob Pike's talk here: https://www.youtube.com/watch?v=HxaD_trXwRE

I've never written a lexer/parser before, and I'm also learning Go
while I do so.  Should be fun :-)

*/
package components

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// StateFunction represents the lexer's current 'state'.  i.e. are we
// at the beginning of an A-instruction, a label or a C-Instruction.
type StateFunction func(*Lexer) StateFunction

// Lexer tracks the progress as the lexer process moves through the
// input string.
type Lexer struct {
	input string         // entire source file, not bothering with streaming (for now)
	start int            // start of current item in bytes, NOT characters
	pos   int            // the position as we search along/end of current item
	width int            // width of last rune that was read
	items chan AsmLexine // channel on which to pass back the tokens
}

////////////////////////////////////////////////////////////////////////////////
// Public

// NewLexer returns both a lexer structure, and its output channel, on
// which 'lexenes' (is that an actual word?) get passed as they are
// read.
func NewLexer(input string) (*Lexer, chan AsmLexine) {
	l := &Lexer{
		input: input,
		items: make(chan AsmLexine)}

	go l.Run()

	return l, l.items
}

// Run starts the lexer process.
func (l *Lexer) Run() {
	state := skipWhitespace(l)

	for state != nil {
		state = state(l)
	}

	close(l.items)
}

////////////////////////////////////////////////////////////////////////////////
// helper routines

// acceptUntil will keep moving through the input string until either
// on of the inValid characters is reached, or white-space or EOL/EOF
// are reached.
func (l *Lexer) acceptUntil(inValid string) {
	for {
		next := l.nextRune()

		isInvalid := strings.IndexRune(inValid, next) >= 0
		isWhiteSpace := isWhiteSpace(next)

		if isInvalid || isWhiteSpace || l.atEof() {
			l.rewind()
			break
		}
	}
}

// ignore moves the start position up to the current 'read' thingie.
// Used if we want to throw away what we've just traversed (comments
// etc).
func (l *Lexer) ignore() {
	l.start = l.pos
}

// emit will thrown the value of pos-start from input onto the output
// channel
func (l *Lexer) emit(aI AsmInstruction) {
	l.items <- AsmLexine{
		Instruction: aI,
		Value:       l.input[l.start:l.pos]}

	l.start = l.pos
	l.width = 0
}

// skipSpaces moves both pos ad start forward until it hits a non
// white-space character.
func (l *Lexer) skipSpaces() {
	for {
		next := string(l.nextRune())

		if next != " " && next != "\t" && next != "\n" {
			l.rewind()
			break
		}

		l.ignore()
	}
}

// isWhiteSpace returns true if the provided rune is a space, tab or
// newline.
func isWhiteSpace(r rune) bool {
	test := string(r)

	return test != " " &&
		test != "\t" &&
		test != "\n"
}

// peek peeks at the next rune, without advancing through the string.
func (l *Lexer) peek() rune {
	result := l.nextRune()
	l.rewind()

	return result
}

// nextRune returns the next rune in input, moving forward.
func (l *Lexer) nextRune() rune {
	rune, width := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += width
	l.width = width

	return rune
}

// error - I should use this method...
func (l *Lexer) error(msg string, args ...interface{}) StateFunction {
	l.items <- AsmLexine{
		Instruction: asmERROR,
		Value:       fmt.Sprintf(msg, args)}

	return nil
}

// atEof ....  Pretty sure you can figure this one out.
func (l *Lexer) atEof() bool {
	return l.start == len(l.input)
}

// rewind moves back to the previous rune.
func (l *Lexer) rewind() {
	l.pos -= l.width
}

// RewindTo rewinds NOTH pos and start to a set location.
func (l *Lexer) RewindTo(p int) error {
	if p < l.start {
		return errors.New("Attempted to rewind past start")
	}

	l.pos = p

	return nil
}

func (l *Lexer) atBeginComment() bool {
	// If the rune at start is '/' and so is the one immediately
	// following, we're at a comment.

	if l.atEof() {
		return false
	}

	initialPos := l.pos
	first := string(l.nextRune())

	if first != "/" || l.atEof() {
		l.rewind()
		return false
	}

	second := string(l.nextRune())
	l.RewindTo(initialPos)

	return second == "/"
}

func (l *Lexer) movePastEol() {
	i := strings.Index(l.input[l.start:], "\n")

	if i == -1 {
		length := len(l.input)
		l.start = length
		l.pos = length

		return
	}

	// make sure that \n isn't the very last character
	if i == len(l.input) {
		i -= 1
	}

	l.start = i + 1
	l.pos = i + 1
}

////////////////////////////////////////////////////////////////////////////////
// State functions

// General "move forward until we find something useful" function.
func skipWhitespace(l *Lexer) StateFunction {
	for {
		l.skipSpaces()

		if !l.atBeginComment() {
			break
		}

		for l.atBeginComment() {
			l.movePastEol()
		}
	}

	if l.atEof() {
		return nil
	}

	next := string(l.peek())

	if next == "@" {
		return aInstruction
	}

	if next == "(" {
		return atLabel
	}

	return cInstruction
}

// Handles A-instructions (@123 /@loop)
func aInstruction(l *Lexer) StateFunction {
	l.acceptUntil("\t\n ")
	l.emit(asmAINSTRUCT)

	return skipWhitespace
}

// Handles labels (i.e. (LOOP))
func atLabel(l *Lexer) StateFunction {
	l.acceptUntil(")\n")
	l.emit(asmLABEL)

	return skipWhitespace
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
func cInstruction(l *Lexer) StateFunction {
	l.acceptUntil("=;")
	next := string(l.peek())

	// destination part of d=c;j
	if next == "=" {
		l.emit(asmDEST)
		return atCmp
	}

	// comp part of c;j
	if next == ";" {
		l.emit(asmCOMP)
		return atJmp
	}

	// comp only (legal, but doesn't really achieve anything)
	l.emit(asmCOMP)
	return skipWhitespace
}

// Handles the 'c' part of a C-Instruction (d=c;j)
func atCmp(l *Lexer) StateFunction {
	l.acceptUntil(";")
	l.emit(asmCOMP)

	return atJmp
}

// Handles the 'j' part of a C-Instruction
func atJmp(l *Lexer) StateFunction {
	l.acceptUntil("")
	l.emit(asmJUMP)

	return skipWhitespace
}
