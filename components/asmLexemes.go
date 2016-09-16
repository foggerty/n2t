package components

import (
	"fmt"
)

type AsmInstruction int

const (
	asmAINSTRUCT AsmInstruction = iota // @123 or @sum
	asmLABEL                           // e.g. (END)
	asmDEST                            // dest part of a c-instruction
	asmCOMP                            // comp part of a c-instruction
	asmJUMP                            // jump part of a c-instruction
	asmEOL                             // end of line, marks end of instruction
	asmEOF                             // end of file
	asmERROR                           // something went horribly wrong
)

type AsmLexeme struct {
	instruction AsmInstruction
	value       string
	lineNum     int
}

func (l AsmLexeme) String() string {
	switch l.instruction {
	case asmEOF:
		return "EOF"
	case asmEOL:
		return fmt.Sprintf("(%d) EOL", l.lineNum)
	case asmAINSTRUCT:
		return fmt.Sprintf("(%d) @%s", l.lineNum, l.value)
	case asmLABEL:
		return fmt.Sprintf("(%d) (%s)", l.lineNum, l.value)
	case asmDEST:
		return fmt.Sprintf("(%d) dst - %s", l.lineNum, l.value)
	case asmJUMP:
		return fmt.Sprintf("(%d) jmp - %s", l.lineNum, l.value)
	case asmCOMP:
		return fmt.Sprintf("(%d) cmp - %s", l.lineNum, l.value)
	case asmERROR:
		return "ERROR - " + l.value
	default:
		panic("Ohshitohshitohshitohshit")
	}
}

type CompilerInstruction int

const (
	cmpThing1 CompilerInstruction = iota
	cmpThing2
)

type CompLexeme struct {
	instruction CompilerInstruction
	value       string
}
