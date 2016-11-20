/*
 State functions and output channel for the assembler's implementation
 of the lexer.
*/

package components

////////////////////////////////////////////////////////////////////////////////
// character sets for various tokens (symbol, instruction etc)
////////////////////////////////////////////////////////////////////////////////

const validSymbol string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.$:-"
const validInstruction string = "-+01!AMD&|nullJGELMNTQEP"

////////////////////////////////////////////////////////////////////////////////
// Here's where it all goes wrong....
////////////////////////////////////////////////////////////////////////////////

func StartLexingAsm(input string) []asmLexeme {
	lex := newLexer(input)

	lex.Run(initState)

	return lex.output
}

////////////////////////////////////////////////////////////////////////////////
// function used by ASM lexer to map current token.
////////////////////////////////////////////////////////////////////////////////

func emit(l *lexer, i asmInstruction) {
	var value string

	switch i {
	case asmEOL:
		fallthrough
	case asmEOF:
		value = ""
	case asmERROR:
		value = l.currentLine()
	default:
		value = l.value()
	}

	lex := asmLexeme{
		lineNum:     l.lineNum,
		instruction: i,
		value:       value,
	}

	l.output = append(l.output, lex)

	if i == asmEOL {
		l.lineNum++
	}
}

////////////////////////////////////////////////////////////////////////////////
// ASM Lexer State Functions
////////////////////////////////////////////////////////////////////////////////

// Skips leading white space, comments and newlines until we reach what is
// hopefully code.
func initState(l *lexer) stateFunction {
	l.skipWhiteSpace()

	if l.atEOF() {
		emit(l, asmEOF)
		return nil
	}

	if l.atEOL() {
		emit(l, asmEOL)
		l.skipEol()
		return initState
	}

	// determine what we're looking at
	next := l.nextInstance("=;(@")

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

	emit(l, asmDEST)

	// move past '='
	l.skipOne()

	return atComp
}

func atComp(l *lexer) stateFunction {
	l.accept(validInstruction)

	if l.nothingFound() {
		return errorState
	}

	emit(l, asmCOMP)

	next := l.peek()

	if next == ";" {
		// move past ';'
		l.skipOne()
		return atJmp
	}

	return endOfInstruction
}

func atAInstruct(l *lexer) stateFunction {
	// move past '@'
	l.skipOne()

	l.accept(validSymbol)

	if l.nothingFound() {
		return errorState
	}

	emit(l, asmAINSTRUCT)

	return endOfInstruction
}

func atJmp(l *lexer) stateFunction {
	l.accept(validInstruction)

	if l.nothingFound() {
		return errorState
	}

	emit(l, asmJUMP)

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

	emit(l, asmLABEL)
	l.skipOne()

	return endOfInstruction
}

func endOfInstruction(l *lexer) stateFunction {
	l.skipWhiteSpace()

	if !(l.atEOL() || l.atEOF()) {
		return errorState
	}

	return initState
}

func errorState(l *lexer) stateFunction {
	emit(l, asmERROR)
	l.skipToEol()

	return initState
}
