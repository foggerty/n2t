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
	Instruction AsmInstruction
	Value       string
	LineNum     int
}

func (l AsmLexeme) String() string {
	switch l.Instruction {
	case asmEOF:
		return "EOF"
	case asmEOL:
		return fmt.Sprintf("(%d) EOL", l.LineNum)
	case asmAINSTRUCT:
		return fmt.Sprintf("(%d) @%s", l.LineNum, l.Value)
	case asmLABEL:
		return fmt.Sprintf("(%d) (%s)", l.LineNum, l.Value)
	case asmDEST:
		return fmt.Sprintf("(%d) dst - %s", l.LineNum, l.Value)
	case asmJUMP:
		return fmt.Sprintf("(%d) jmp - %s", l.LineNum, l.Value)
	case asmCOMP:
		return fmt.Sprintf("(%d) cmp - %s", l.LineNum, l.Value)
	case asmERROR:
		return "ERROR - " + l.Value
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
