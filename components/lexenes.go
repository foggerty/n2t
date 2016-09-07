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

type AsmLexine struct {
	Instruction AsmInstruction
	Value       string
}

func (l AsmLexine) String() string {
	switch l.Instruction {
	case asmEOF:
		return "EOF"
	case asmAINSTRUCT:
		return fmt.Sprintf("@%s", l.Value)
	case asmLABEL:
		return fmt.Sprintf("(%s)", l.Value)
	case asmDEST:
		return fmt.Sprintf("dst - %s", l.Value)
	case asmJUMP:
		return fmt.Sprintf("jmp - %s", l.Value)
	case asmCOMP:
		return fmt.Sprintf("cmp - %s", l.Value)
	default:
		return "I have no idea."
	}
}

type CompilerInstruction int

const (
	cmpThing1 CompilerInstruction = iota
	cmpThing2
)

type CompLexene struct {
	instruction CompilerInstruction
	value       string
}
