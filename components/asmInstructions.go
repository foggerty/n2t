package components

////////////////////////////////////////////////////////////////////////////////
// All instructions are 16 bits.  Signing doesn't really come into play here,
// as any 'numbers' are encoded in the lower 15 bits (although they are
// using 2's compliment).

type asm uint16

const cInst asm = 7 << 13
const aInst asm = 0

////////////////////////////////////////////////////////////////////////////////
// Register constants - these are all RAM addresses.

var registers = map[string]asm{
	"R0":  0,
	"R1":  1,
	"R2":  2,
	"R3":  3,
	"R4":  4,
	"R5":  5,
	"R6":  6,
	"R7":  7,
	"R8":  8,
	"R9":  9,
	"R10": 10,
	"R11": 11,
	"R12": 12,
	"R13": 13,
	"R14": 14,
	"R15": 15,
}

var pointers = map[string]asm{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"SCREEN": 16384,
	"KBD":    24576,
}

////////////////////////////////////////////////////////////////////////////////
// Destination part of instruction.
// Shift three bits to the left for their position in an instruction.

var destMap = map[string]asm{
	"null": 0,
	"M":    1 << 3,
	"D":    2 << 3,
	"MD":   3 << 3,
	"A":    4 << 3,
	"AM":   5 << 3,
	"AD":   6 << 3,
	"AMD":  7 << 3,
}

////////////////////////////////////////////////////////////////////////////////
// Jump part of instruction.
// No shifting required, they're the three LSBs.

var jmpMap = map[string]asm{
	"null": 0,
	"JGT":  1,
	"JEQ":  2,
	"JGE":  3,
	"JLT":  4,
	"JNE":  5,
	"JLE":  6,
	"JMP":  7,
}

////////////////////////////////////////////////////////////////////////////////
// Compute part of instruction.

var cmpMap = map[string]asm{
	// A = 0
	"0":   42 << 6,
	"1":   63 << 6,
	"-1":  58 << 6,
	"D":   12 << 6,
	"A":   48 << 6,
	"!D":  13 << 6,
	"!A":  49 << 6,
	"-D":  15 << 6,
	"-A":  51 << 6,
	"D+1": 31 << 6,
	"A+1": 55 << 6,
	"D-1": 14 << 6,
	"A-1": 50 << 6,
	"D+A": 2 << 6,
	"D-A": 19 << 6,
	"A-D": 7 << 6,
	"D&A": 0 << 6,
	"D|A": 21 << 6,

	// A = 1
	"M":   112 << 6,
	"!M":  113 << 6,
	"-M":  79 << 6,
	"M+1": 119 << 6,
	"M-1": 114 << 6,
	"D+M": 66 << 6,
	"D-M": 83 << 6,
	"M-D": 71 << 6,
	"D&M": 64 << 6,
	"D|M": 85 << 6,
}
