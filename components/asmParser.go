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

const maxConst = 32768 // 2^15

func newParser(input chan asmLexeme) (*asmParser, errorList) {
	parser := asmParser{
		items:       input,
		output:      make(chan asm),
		symbolTable: newSymbolTable(),
	}

	// first pass, building symbol table and recording errors
	errs := parser.buildSymbols()
	parser.errored = errs != nil
	parser.witeMem()

	return &parser, errs
}

func (p *asmParser) run(f *os.File) errorList {
	var errs []error

	if p.errored {
		return append(errs, errors.New("Cannot parse - errors found by lexer"))
	}

	var i asm
	var err error

	for _, lex := range p.lexemes {

		switch lex.instruction {
		case asmAINSTRUCT:
			i, err = p.mapToA(lex)
		default:
			continue
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}

		fmt.Fprintf(f, "%.16b\n", i)
	}

	return errs
}

func (p *asmParser) mapToA(l asmLexeme) (asm, error) {
	// is it a constant?
	if c, err := strconv.Atoi(l.value); err == nil {
		// and is it within the allowed range? (0 - 2^15)
		if c >= 0 && c <= maxConst {
			return aInst | asm(c), nil
		}

		return 0, fmt.Errorf("Constant value out of range: %s", l.value)
	}

	// does the value exist in the symbol table?
	if sym, err := p.symbolValue(l.value); err == nil {
		return aInst | sym, nil
	}

	// is it a register?
	if reg, ok := registers[l.value]; ok {
		return aInst | reg, nil
	}

	// is it a predefined pointer?
	if ptr, ok := pointers[l.value]; ok {
		return aInst | ptr, nil
	}

	return asm(0), fmt.Errorf("Unrecognised value for A-Instruction: %s", l.value)
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
			p.addLabel(lex.value, asm(line-1))

		case asmAINSTRUCT:
			if !isInt(lex.value) ||
				!isRegister(lex.value) {
				p.addVariable(lex.value)
			}
		}
	}

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
