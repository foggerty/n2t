package components

import (
	"strings"
	"unicode/utf8"
)

// StateFunction represents the lexer's current 'state'.  i.e. are we
// at the beginning of an A-instruction, a label or a C-Instruction.
type stateFunction func(*lexer) stateFunction

// Lexer tracks the progress as the lexer process moves through the
// input string.  Todo - 'instruct' should be inferred by the current
// state function somehow and maybe passed in to the emit function,
// this is brittle code.
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

	state := initState

	for state != nil {
		state = state(l)
	}
}

////////////////////////////////////////////////////////////////////////////////
// character sets for various tokens (symbol, instruction etc)

const validSymbol string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.$:-"
const validInstruction string = "-+01!AMD&|nullJGELMNTQEP"

////////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////////

// True if at EOF
func (l *lexer) atEOF() bool {
	return l.pos >= len(l.input) || l.input == ""
}

func (l *lexer) atEOL() bool {
	next := l.peek()

	return next == "\r" || next == "\n"
}

// True if sitting at the beginning of a comment.
func (l *lexer) atComment() bool {
	if l.pos < len(l.input)-2 {
		return "//" == l.input[l.pos:l.pos+2]
	}

	return false
}

// Starting from current position, look for the next ';', '=', '@' or
// '(' up to either EOL or EOF.  If there isn't one, returns empty
// string.
func (l *lexer) nextSeparator() string {
	i := strings.IndexAny(l.input[l.pos:], ";=(@")

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

// Puts current selection onto the output channel.
// Note that ONLY state functions should be calling this.
func (l *lexer) emit(i asmInstruction) {
	var value string

	switch i {
	case asmEOL:
		fallthrough
	case asmEOF:
		value = ""
	case asmERROR:
		value = l.currentLine()
	default:
		value = l.input[l.start:l.pos]
	}

	lex := asmLexeme{
		lineNum:     l.lineNum,
		instruction: i,
		value:       value,
	}

	l.items <- lex

	if i == asmEOL {
		l.lineNum++
	}
}

// Skips over white space and comment until EOL or EOF.
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

////////////////////////////////////////////////////////////////////////////////
// State Functions
////////////////////////////////////////////////////////////////////////////////

// Skips leading white space, comments and newlines until we reach what is
// hopefully code.
func initState(l *lexer) stateFunction {
	l.skipWhiteSpace()

	if l.atEOF() {
		l.emit(asmEOF)
		return nil
	}

	if l.atEOL() {
		l.emit(asmEOL)
		l.skipEol()
		return initState
	}

	// determine what we're looking at
	next := l.nextSeparator()

	switch next {
	case "=":
		// instruction has a DEST part
		return atDest
	case ";":
		// only COMP & JMP
		return atComp
	case "(":
		// at a label
		return atLabel
	case "@":
		return atAInstruct
	default:
		// anything else must be an error
		return errorState
	}
}

func atDest(l *lexer) stateFunction {
	l.accept(validInstruction)
	next := l.peek()

	if l.nothingFound() || next != "=" {
		return errorState
	}

	l.emit(asmDEST)

	// move past '='
	l.skipOne()

	return atComp
}

func atComp(l *lexer) stateFunction {
	l.accept(validInstruction)

	if l.nothingFound() {
		return errorState
	}

	l.emit(asmCOMP)

	next := l.peek()

	if next == ";" {
		// move past ';'
		l.skipOne()
		return atJmp
	}

	return endOfInstruction
}

func atAInstruct(l *lexer) stateFunction {
	_ = "breakpoint"
	// move past '@'
	l.skipOne()

	l.accept(validSymbol)

	if l.nothingFound() {
		return errorState
	}

	l.emit(asmAINSTRUCT)

	return endOfInstruction
}

func atJmp(l *lexer) stateFunction {
	l.accept(validInstruction)

	if l.nothingFound() {
		return errorState
	}

	l.emit(asmJUMP)

	return endOfInstruction
}

func atLabel(l *lexer) stateFunction {
	// move over '('
	l.skipOne()
	l.accept(validSymbol)
	next := l.peek()

	if l.nothingFound() || next != ")" {
		return errorState
	}

	l.emit(asmLABEL)
	l.skipOne()

	return endOfInstruction
}

func endOfInstruction(l *lexer) stateFunction {
	// remove this state and your code can have multiple instructions
	// per line, but trying to stick to the 'spec' :-)
	l.skipWhiteSpace()

	if !(l.atEOL() || l.atEOF()) {
		return errorState
	}

	return initState
}

func errorState(l *lexer) stateFunction {
	l.emit(asmERROR)
	l.skipToEol()

	return initState
}
