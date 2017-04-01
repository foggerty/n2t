package components

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// AsmParser represents a n2t parser for the assembler.  It takes in a
// channel of lexemes, and provided they look ok, will then start
// passing instructions (as strings) back via the output channel.
type AsmParser struct {
	items  chan asmLexeme
	Output chan string
	symbolTable
	lexemes []asmLexeme
	Error   error
}

const maxConst = 32768 // 2^15

// NewParser creates a new instance of AsmParser, kicks off the
// process by running the first pass (to build symbol table) and
// returns the parser.  Any errors encountered during the first
// pass will be attached to the Error field.
func NewParser(input chan asmLexeme) AsmParser {
	parser := AsmParser{
		items:       input,
		Output:      make(chan string),
		symbolTable: newSymbolTable(),
	}

	// first pass, building symbol table and recording errors
	parser.buildSymbols()

	if parser.Error == nil {
		go parser.run()
	}

	return parser
}

// Run is badly named and is about to be changed.  Note that the Lexer
// is actually doing a lot of error checking, so can assume at this
// point that, while that may they not be correctly spelled, we're not
// going to get more than one jmp per line etc, or more than three
// parts (d=c;j) per line.  So this isn't really a parser, it just
// maps instruction mnemonics.
//
// At the first error will stop writing to the output channel, but
// sill continue to parse the rest of the lexemes, so that a full list
// of errors can still be returned.
func (p *AsmParser) run() {
	defer close(p.Output)

	var errs errorList

	if p.Error != nil {
		return
	}

	var i asm // instruction, reset to 0 after every write
	var err error
	var d, c, j asm // dest, comp, jump, OR together for final instruction

	writeResult := func() {
		if err != nil {
			errs = append(errs, err)
		}

		if errs == nil {
			p.Output <- fmt.Sprintf("%.16b", i)
		}

		i = 0
	}

	for index, lex := range p.lexemes {

		switch lex.instruction {

		// possible edge case, hitting EOF before an EOL
		case asmEOF:
			fallthrough

		case asmEOL:
			prev := p.previousInstruction(index)

			if prev.instruction != asmLABEL {
				writeResult()
			}

		case asmAINSTRUCT:
			prev := p.previousInstruction(index)

			if prev.instruction == asmAINSTRUCT {
				fmt.Fprintf(os.Stderr, "WARNING - redundant loading of A-Register on line %d\n", prev.lineNum)
			}

			i, err = p.mapToA(lex)

		case asmLABEL:
			index += 2 // skip label and EOL
			continue

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

	p.Error = errs.asError()
}

func (p *AsmParser) previousInstruction(index int) asmLexeme {
	nil := asmLexeme{instruction: asmNULL}

	if index-2 < 0 {
		return nil
	}

	previous := p.lexemes[index-1]

	if previous.instruction == asmEOL {
		previous = p.lexemes[index-2]
	}

	return previous
}

func (p *AsmParser) mapToA(l asmLexeme) (asm, error) {
	// is it a constant?
	if c, err := strconv.Atoi(l.value); err == nil {
		// and is it within the allowed range? (0 - 2^15)
		if c >= 0 && c <= maxConst {
			return aInst | asm(c), nil
		}

		return 0, fmt.Errorf("Constant value out of range, line %d: %s", l.lineNum, l.value)
	}

	// does the value exist in the symbol table?
	if sym := p.symbolValue(l.value); sym != asm(0) {
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

// First pass (parse?) - builds the symbol table.
func (p *AsmParser) buildSymbols() {
	var pCount = 0 // instruction memory counter
	var errs errorList
	var foundComp bool
	var previous = asmEOL

	for {
		lex, ok := <-p.items

		if !ok {
			break
		}

		switch lex.instruction {

		case asmERROR:
			errs = append(errs, errors.New(lex.value))

		case asmEOL:
			if foundComp {
				pCount++
				foundComp = false
			}

		// case asmEOF:
		// 	break

		case asmLABEL:
			p.addLabel(lex.value, asm(pCount))

		case asmAINSTRUCT:
			pCount++
			if !isInt(lex.value) ||
				!isRegister(lex.value) {
				p.addVariable(lex.value)
			}

		case asmDEST:
			fallthrough
		case asmCOMP:
			foundComp = true
		}

		dupeEol := lex.instruction == asmEOL && previous == asmEOL

		if !dupeEol {
			p.lexemes = append(p.lexemes, lex)
		}

		previous = lex.instruction
	}

	if len(errs) == 0 {
		p.writeMem()
	}

	p.Error = errs.asError()
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
