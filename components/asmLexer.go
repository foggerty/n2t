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

	lex := newLexer(input)

	lex.Run(initState)

	return lex.output
}

////////////////////////////////////////////////////////////////////////////////
// ASM Lexer State Functions
////////////////////////////////////////////////////////////////////////////////

// Skips leading white space, comments and newlines until we reach what is
// hopefully code.
func initState(l* lexer) stateFunction {

	l.skipWhiteSpace()

	if l.atEOF() {
		l.emit(asmEOF)
		close(l.output)
		return nil
	}

	if l.atEOL() {
		l.emit(asmEOL)
		l.skipEol()
		return initState(l)
	}

	// determine what we're looking at
	next := l.nextInstance("=;(@")

	switch next {
	case "=":
		// instruction has a DEST part
		return atDest(l)
	case ";":
		// only COMP & JMP
		return atComp(l)
	case "(":
		// at a label
		return atLabel(l)
	case "@":
		return atAInstruct(l)
	default:
		// anything else must be an error
		return errorState(l)
	}
}

func atDest(l* lexer) stateFunction {

	l.accept(validInstruction)
	next := l.peek()

	if l.nothingFound() || next != "=" {
		return errorState(l)
	}

	l.emit(asmDEST)

	// move past '='
	l.skipOne()

	return atComp(l)
}

func atComp(l* lexer) stateFunction {

	l.accept(validInstruction)

	if l.nothingFound() {
		return errorState(l)
	}

	l.emit(asmCOMP)

	next := l.peek()

	if next == ";" {
		// move past ';'
		l.skipOne()
		return atJmp(l)
	}

	return endOfInstruction(l)
}

func atAInstruct(l* lexer) stateFunction {

	// move past '@'
	l.skipOne()

	l.accept(validSymbol)

	if l.nothingFound() {
		return errorState(l)
	}

	l.emit(asmAINSTRUCT)

	return endOfInstruction(l)
}

func atJmp(l* lexer) stateFunction {

	l.accept(validInstruction)

	if l.nothingFound() {
		return errorState(l)
	}

	l.emit(asmJUMP)

	return endOfInstruction(l)
}

func atLabel(l* lexer) stateFunction {

	// move over '('
	l.skipOne()
	l.accept(validSymbol)
	next := l.peek()

	if l.nothingFound() || next != ")" {
		return errorState(l)
	}

	l.emit(asmLABEL)
	l.skipOne()

	return endOfInstruction(l)
}

func endOfInstruction(l* lexer) stateFunction {

	l.skipWhiteSpace()

	if !(l.atEOL() || l.atEOF()) {
		return errorState(l)
	}

	return initState(l)
}

func errorState(l* lexer) stateFunction {

	l.emit(asmERROR)
	l.skipToEol()

	return initState(l)
}
