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

func StartLexingAsm(input string) chan asmLexeme {
	out := make(chan asmLexeme)

	lex := newLexer(input)

	lex.Run(initState(out))

	return out
}

////////////////////////////////////////////////////////////////////////////////
// function used by ASM lexer to map current token.
////////////////////////////////////////////////////////////////////////////////

func emit(l *lexer, i asmInstruction, c chan asmLexeme) {
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

	c <- lex

	if i == asmEOL {
		l.lineNum++
	}
}

////////////////////////////////////////////////////////////////////////////////
// ASM Lexer State Functions
////////////////////////////////////////////////////////////////////////////////

// Skips leading white space, comments and newlines until we reach what is
// hopefully code.
func initState(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.skipWhiteSpace()

		if l.atEOF() {
			emit(l, asmEOF, c)
			close(c)
			return nil
		}

		if l.atEOL() {
			emit(l, asmEOL, c)
			l.skipEol()
			return initState(c)
		}

		// determine what we're looking at
		next := l.nextInstance("=;(@")

		switch next {
		case "=":
			// instruction has a DEST part
			return atDest(c)
		case ";":
			// only COMP & JMP
			return atComp(c)
		case "(":
			// at a label
			return atLabel(c)
		case "@":
			return atAInstruct(c)
		default:
			// anything else must be an error
			return errorState(c)
		}
	}

}

func atDest(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.accept(validInstruction)
		next := l.peek()

		if l.nothingFound() || next != "=" {
			return errorState(c)
		}

		emit(l, asmDEST, c)

		// move past '='
		l.skipOne()

		return atComp(c)
	}

}

func atComp(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.accept(validInstruction)

		if l.nothingFound() {
			return errorState(c)
		}

		emit(l, asmCOMP, c)

		next := l.peek()

		if next == ";" {
			// move past ';'
			l.skipOne()
			return atJmp(c)
		}

		return endOfInstruction(c)
	}

}

func atAInstruct(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		// move past '@'
		l.skipOne()

		l.accept(validSymbol)

		if l.nothingFound() {
			return errorState(c)
		}

		emit(l, asmAINSTRUCT, c)

		return endOfInstruction(c)
	}

}

func atJmp(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.accept(validInstruction)

		if l.nothingFound() {
			return errorState(c)
		}

		emit(l, asmJUMP, c)

		return endOfInstruction(c)
	}

}

func atLabel(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		// move over '('
		l.skipOne()
		l.accept(validSymbol)
		next := l.peek()

		if l.nothingFound() || next != ")" {
			return errorState(c)
		}

		emit(l, asmLABEL, c)
		l.skipOne()

		return endOfInstruction(c)
	}

}

func endOfInstruction(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.skipWhiteSpace()

		if !(l.atEOL() || l.atEOF()) {
			return errorState(c)
		}

		return initState(c)
	}

}

func errorState(c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		emit(l, asmERROR, c)
		l.skipToEol()

		return initState(c)
	}

}
