package components

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type asmParser struct {
	items  chan asmLexeme
	output chan asm
	symbolTable
	lexemes []asmLexeme
	errored bool
}

func newParser(input chan asmLexeme) (*asmParser, errorList) {
	parser := asmParser{
		items:       input,
		output:      make(chan asm),
		symbolTable: newSymbolTable(),
	}

	// first pass, building symbol table and recording errors
	errs := parser.buildSymbols()

	return &parser, errs
}

func (p *asmParser) run(f *os.File) errorList {
	// if errored, we don't write to the file, but do parse the lexemes
	// looking for additional errors.

	for lex := range p.lexemes {
		fmt.Println(lex)
	}

	return nil
}

func (p *asmParser) buildSymbols() errorList {
	// To-do: find out what the actual instruction memory offset is
	line := 1
	var errs []error

	for lex := range p.items {
		// will need these for the second pass
		p.lexemes = append(p.lexemes, lex)

		switch lex.instruction {

		case asmERROR:
			msg := fmt.Sprintf("%q", lex)
			errs = append(errs, errors.New(msg))

		case asmEOL:
			// the lexer will compact extra EOL chars, so there will be a
			// 1-1 relationship between EOL count and instruction memory offset
			line++

		case asmEOF:
			break

		case asmLABEL:
			p.addLabel(lex.value, asm(line+1))

		case asmAINSTRUCT:
			if !isInt(lex.value) ||
				!isRegister(lex.value) {
				p.addVariable(lex.value)
			}
		}
	}

	p.errored = errs != nil

	return errs
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)

	return err != nil
}

func isRegister(s string) bool {
	_, ok := registers[s]

	return ok
}
