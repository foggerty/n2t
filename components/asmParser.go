/*

  asmParser

  Input: a channel of asmLexemes

  Output: a channel of asm (i.e. 16 bit words, each representing a
  single instruction).

*/

package components

import (
	"strconv"
)

type asmParser struct {
	items  chan AsmLexeme
	output chan asm
	symbolTable
	lexemes []AsmLexeme
}

func newParser(i chan AsmLexeme) <-chan asm {
	p := asmParser{
		items:       i,
		output:      make(chan asm),
		symbolTable: newSymbolTable(),
	}

	p.run()

	return p.output
}

func (p *asmParser) run() {
	// collect all lexemes and build the symbol table (first pass)
	p.buildSymbols()

	// populate memory locations in symbol table
	p.initMemory()

	// second pass, turn into code into an output channel of ints
}

func (p *asmParser) buildSymbols() {
	// To-do: find out what the actual memory offset is
	line := 1

	for lex := range p.items {
		// will need these for the second pass
		p.lexemes = append(p.lexemes, lex)

		switch lex.instruction {

		case asmEOL:
			line++

		case asmEOF:
			break

		case asmLABEL:
			p.addLabel(lex.value, asm(line))

		case asmAINSTRUCT:
			if !isInt(lex.value) ||
				!isRegister(lex.value) {
				p.addVariable(lex.value)
			}
		}
	}
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)

	return err != nil
}

func isRegister(s string) bool {
	_, ok := registers[s]

	return ok
}
