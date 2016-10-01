package components

import (
	"strings"
	"unicode/utf8"
)

// StateFunction represents the lexer's current 'state', and hopefully
// knows what to do from this point forward.  Returns nil if there's
// nothing left to do.
type stateFunction func(*lexer) stateFunction

// Lexer tracks the progress as the lexer process moves through the
// input string.
type lexer struct {
	input    string // entire source file, not bothering with streaming (for now)
	start    int    // start of current item in bytes, NOT characters
	pos      int    // the position as we search along/end of current item
	width    int    // width of last rune that was read
	lineNum  int    // current source line number
	instruct bool   // true if we've started processing an instruction
}

// Requires an initial state function to run.
func newLexer(input string) *lexer {
	return &lexer{
		input:   input,
		lineNum: 1,
	}
}

// Kick off the lexing process.
func (l *lexer) Run(init stateFunction) {
	state := init

	go func() {
		for state != nil {
			state = state(l)
		}
	}()
}

// Returns the current token.
func (l *lexer) value() string {
	return l.input[l.start:l.pos]
}

// True if at EOF
func (l *lexer) atEOF() bool {
	return l.pos >= len(l.input) || l.input == ""
}

// True id at EOL
func (l *lexer) atEOL() bool {
	next := l.peek()

	return next == "\r" || next == "\n"
}

// True if sitting at the beginning of a comment.
// To do - consumer provides function to determine this.
func (l *lexer) atComment() bool {
	if l.pos < len(l.input)-2 {
		return "//" == l.input[l.pos:l.pos+2]
	}

	return false
}

// Starting from current position, look for the next ';', '=', '@' or
// '(' up to either EOL or EOF.  If there isn't one, returns empty
// string.
func (l *lexer) nextInstance(of string) string {
	i := strings.IndexAny(l.input[l.pos:], of)

	if i == -1 {
		return ""
	}

	i += l.pos

	// test to see if separator is on the following line
	iEol := l.nextEol() + l.pos

	if iEol-i < 0 {
		return ""
	}

	return l.input[i : i+1]
}

// Return the entire line (i.e. back from start to the next EOL/BOF,
// forward toward the next EOL/EOF).
func (l *lexer) currentLine() string {
	var start int
	end := l.nextEol()

	if l.lineNum == 1 {
		start = 0
	} else {
		// For linenum to be anything but 1, a EOL has to have been
		// processed.
		start = l.start - 1
		for {
			test := l.input[start : start+1]
			if test == "\n" {
				start++
				break
			}
			start--
		}
	}

	return l.input[start:end]
}

// Returns index of next EOL character, or EOF position if none found
func (l *lexer) nextEol() int {
	i := strings.IndexAny(l.input[l.pos:], "\r\n")

	if i == -1 {
		return len(l.input)
	}

	return l.pos + i
}

// Skips over white space an possible comment until EOL or EOF.
func (l *lexer) skipWhiteSpace() {
	for {
		l.skipChars(" \t")

		if l.atComment() {
			l.skipToEol()
			continue
		}

		break
	}
}

// Moves both start and pos forward until EOL or EOF.
func (l *lexer) skipToEol() {
	// position of next EOL
	i := l.nextEol()
	l.pos = i
	l.start = i
}

// Skips past characters, moving pos and start forward.
func (l *lexer) skipChars(chars string) {
	for {
		next := l.next()

		if !strings.ContainsAny(next, chars) || next == "" {
			l.rewind()
			break
		}

		l.ignore()
	}
}

// Assumes at EOL
func (l *lexer) skipEol() {
	next := l.next()

	if next == "\r" {
		l.skipOne()
	}

	l.ignore()
}

func (l *lexer) skipOne() {
	l.next()
	l.ignore()
}

// Reads and returns next rune as a string.
func (l *lexer) next() string {
	next, width := utf8.DecodeRuneInString(l.input[l.pos:])

	l.pos += width
	l.width = width

	return string(next)
}

// Rewinds pos back by the width of the last rune read.
func (l *lexer) rewind() {
	l.pos -= l.width
	l.width = 0
}

// Moves pos up to start, ignoring runes in between.
func (l *lexer) ignore() {
	l.start = l.pos
}

// Moves pos forward over accepted characters.
func (l *lexer) accept(chars string) {
	next := l.next()

	for {
		if !strings.ContainsAny(chars, next) {
			break
		}

		next = l.next()
	}

	l.rewind()
}

// True if pos has not advanced past start.
func (l *lexer) nothingFound() bool {
	return l.pos-l.start <= 0
}

// Returns next rune as a sting, or empty string if EOF
func (l *lexer) peek() string {
	if l.atEOF() {
		return ""
	}

	next := l.next()
	l.rewind()

	return string(next)
}
