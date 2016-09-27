package components

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type asmParser struct {
	items chan asmLexeme
	symbolTable
	lexemes []asmLexeme
	errored bool
}

const maxConst = 32768 // 2^15

func newParser(input chan asmLexeme) (*asmParser, errorList) {
	parser := asmParser{
		items:       input,
		symbolTable: newSymbolTable(),
	}

	// first pass, building symbol table and recording errors
	errs := parser.buildSymbols()
	parser.errored = errs != nil
	parser.witeMem()

	return &parser, errs
}

// The Lexer is actually doing a lot of error checking, so can assume
// at this point that, while that may they not be correctly spelled,
// we're not going to get more than one jmp per line etc, or more than
// three parts (d=c;j) per line.  So this isn't really a parser, it
// just maps instruction mnemonics.
func (p *asmParser) run(f *os.File) errorList {
	var errs []error

	if p.errored {
		return errorList{errors.New("Cannot parse - errors found by lexer")}
	}

	var i asm // instruction, reset to 0 after every write
	var err error
	var d, c, j asm // dest, comp, jump
	var index = 0

Loop:
	for {

		lex := p.lexemes[index]

		switch lex.instruction {

		case asmEOF:
			break Loop

		case asmEOL:
			if err == nil {
				fmt.Fprintf(f, "%.16b\n", i)
			} else {
				errs = append(errs, err)
			}

			i = 0

			// A - Instructions

		case asmAINSTRUCT:
			i, err = p.mapToA(lex)

		case asmLABEL:
			index += 2 // skip label and EOL
			continue

			// C - Instructions

		case asmJUMP:
			j, err = mapJmp(lex.value)
			i = i | j

		case asmCOMP:
			c, err = mapCmp(lex.value)
			i = i | c

		case asmDEST:
			d, err = mapDest(lex.value)
			i = i | d
		}

		index++
	}

	return errs
}

func printInstruction(i asm, err error, f *os.File) {
	if err == nil {
		fmt.Fprintf(f, "%.16b\n", i)
	}
}

func (p *asmParser) mapToA(l asmLexeme) (asm, error) {
	// is it a constant?
	if c, err := strconv.Atoi(l.value); err == nil {
		// and is it within the allowed range? (0 - 2^15)
		if c >= 0 && c <= maxConst {
			return aInst | asm(c), nil
		}

		return 0, fmt.Errorf("Constant value out of range, line %d: %s", l.lineNum, l.value)
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
	line := 0
	var errs []error
	var foundComp bool

	for {
		lex, ok := <-p.items

		if !ok {
			break
		}

		// will need these for the second pass
		p.lexemes = append(p.lexemes, lex)

		switch lex.instruction {

		case asmERROR:
			errs = append(errs, errors.New(lex.value))

		case asmEOL:
			if foundComp {
				line++
				foundComp = false
			}

		case asmEOF:
			break

		case asmLABEL:
			p.addLabel(lex.value, asm(line))

		case asmAINSTRUCT:
			line++
			if !isInt(lex.value) ||
				!isRegister(lex.value) {
				p.addVariable(lex.value)
			}

		case asmDEST:
			fallthrough
		case asmCOMP:
			foundComp = true
		}
	}

	return errs
}

func mapInstruction(i string, m map[string]asm) (asm, error) {
	res, ok := m[i]

	if !ok {
		return 0, fmt.Errorf("Unrecognised instruction: %s", i)
	}

	return res | cInst, nil
}

func mapJmp(j string) (asm, error) {
	return mapInstruction(j, jmpMap)
}

func mapDest(d string) (asm, error) {
	return mapInstruction(d, destMap)
}

func mapCmp(c string) (asm, error) {
	return mapInstruction(c, cmpMap)
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)

	return err != nil
}

func isRegister(s string) bool {
	_, r := registers[s]
	_, p := pointers[s]

	return r || p
}
