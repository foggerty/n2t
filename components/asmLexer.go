package components

////////////////////////////////////////////////////////////////////////////////
// character sets for various tokens (symbol, instruction etc)
////////////////////////////////////////////////////////////////////////////////

const validSymbol string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.$:-"
const validInstruction string = "-+01!AMD&|nullJGELMNTQEP"

////////////////////////////////////////////////////////////////////////////////
// Here's where it all goes wrong....
////////////////////////////////////////////////////////////////////////////////

func kickOff(input string) chan asmLexeme {
	out := make(chan asmLexeme)

	lex := newLexer(input)

	lex.Run(initState(lex, out))

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
func initState(l *lexer, c chan asmLexeme) stateFunction {

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
			return initState(l, c)
		}

		// determine what we're looking at
		next := l.nextInstance("=;(@")

		switch next {
		case "=":
			// instruction has a DEST part
			return atDest(l, c)
		case ";":
			// only COMP & JMP
			return atComp(l, c)
		case "(":
			// at a label
			return atLabel(l, c)
		case "@":
			return atAInstruct(l, c)
		default:
			// anything else must be an error
			return errorState(l, c)
		}
	}

}

func atDest(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.accept(validInstruction)
		next := l.peek()

		if l.nothingFound() || next != "=" {
			return errorState(l, c)
		}

		emit(l, asmDEST, c)

		// move past '='
		l.skipOne()

		return atComp(l, c)
	}

}

func atComp(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.accept(validInstruction)

		if l.nothingFound() {
			return errorState(l, c)
		}

		emit(l, asmCOMP, c)

		next := l.peek()

		if next == ";" {
			// move past ';'
			l.skipOne()
			return atJmp(l, c)
		}

		return endOfInstruction(l, c)
	}

}

func atAInstruct(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		// move past '@'
		l.skipOne()

		l.accept(validSymbol)

		if l.nothingFound() {
			return errorState(l, c)
		}

		emit(l, asmAINSTRUCT, c)

		return endOfInstruction(l, c)
	}

}

func atJmp(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.accept(validInstruction)

		if l.nothingFound() {
			return errorState(l, c)
		}

		emit(l, asmJUMP, c)

		return endOfInstruction(l, c)
	}

}

func atLabel(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		// move over '('
		l.skipOne()
		l.accept(validSymbol)
		next := l.peek()

		if l.nothingFound() || next != ")" {
			return errorState(l, c)
		}

		emit(l, asmLABEL, c)
		l.skipOne()

		return endOfInstruction(l, c)
	}

}

func endOfInstruction(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		l.skipWhiteSpace()

		if !(l.atEOL() || l.atEOF()) {
			return errorState(l, c)
		}

		return initState(l, c)
	}

}

func errorState(l *lexer, c chan asmLexeme) stateFunction {

	return func(l *lexer) stateFunction {
		emit(l, asmERROR, c)
		l.skipToEol()

		return initState(l, c)
	}

}
