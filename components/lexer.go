package components

// Taken form Rob Pike's talk here: https://www.youtube.com/watch?v=HxaD_trXwRE

// To-do: how much of this can be factored out, so that it can be
// reused with the compiler?  The run loop is obvious, as well as
// EOF/EOL tests etc.  But can I maybe update it to take a lookup
// table of strings => types?  For reserved words, yes, but variables
// etc, arithmetic statements, probably not.

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type StateFunction func(*Lexer) StateFunction

type Lexer struct {
	input string         // entire source file, not bothering with streaming (for now)
	start int            // start of current item in bytes, NOT characters
	pos   int            // the position as we search along/end of current item
	width int            // number of runes between pos-start
	items chan AsmLexine // channel on which to pass back the tokens
}

func NewLexer(input string) (*Lexer, chan AsmLexine) {
	l := &Lexer{
		input: input,
		items: make(chan AsmLexine)}

	go l.Run()

	return l, l.items
}

func (l *Lexer) Run() {
	state := l.skipWhitespace()

	for state != nil {
		state = state(l)
	}

	close(l.items)
}

// only run when we've gotten to the end of a recognisable token.
func (l *Lexer) emit(t AsmInstruction) {
	l.items <- AsmLexine{
		Instruction: t,
		Value:       l.input[l.start:l.pos]}

	l.start = l.pos
}

func (l *Lexer) startOfToken() StateFunction {
	// get to end of next token...
	// figure out what we got?
	// emit
	return nil
}

// Skips white space and comments
func (l *Lexer) skipWhitespace() StateFunction {
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

	return nil
}

func (l *Lexer) skipSpaces() {
	for {
		next := string(l.nextRune())

		if next != " " && next != "\t" && next != "\n" {
			l.rewind()
			break
		}

		l.start = l.pos
	}
}

// Returns the next rune, and moves pos forward over it
func (l *Lexer) nextRune() rune {
	next, width := utf8.DecodeRuneInString(l.input[l.start:])
	l.pos += width

	if next == utf8.RuneError {
		msg := fmt.Sprintf("Unrecognisable UTF8 char at byte %d", l.start)
		panic(msg)
	}

	return next
}

func (l *Lexer) atEof() bool {
	return l.start == len(l.input)
}

func (l *Lexer) rewind() {
	l.pos = l.start
}

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
		l.RewindTo(initialPos)
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

	l.start = i + 1
	l.pos = i + 1
}
