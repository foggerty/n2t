package components

////////////////////////////////////////////////////////////////////////////////
// All instructions are 16 bits.  Signing doesn't really come into play here,
// as any 'numbers' are encoded in the lower 15 bits (although they are
// using 2's compliment).

type asm uint16

const cInst asm = 7 << 13
const aInst asm = 0

////////////////////////////////////////////////////////////////////////////////
// Register constants

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
// Literals and map for destination part of instruction.
// Need to be shifted three to the left for their position in an instruction

const (
	dstNull asm = 0
	dstM        = 1
	dstD        = 2
	dstMD       = 3
	dstA        = 4
	dstAM       = 5
	dstAD       = 6
	dstAMD      = 7
)

var destMap = map[string]asm{
	"null": dstNull << 3,
	"M":    dstM << 3,
	"D":    dstD << 3,
	"MD":   dstMD << 3,
	"A":    dstA << 3,
	"AM":   dstAM << 3,
	"AD":   dstAD << 3,
	"AMD":  dstAMD << 3,
}

////////////////////////////////////////////////////////////////////////////////
// Literals and map for jump part of instruction.
// No shifting required, they're the three LSBs.

const (
	jmpNull asm = 0
	jmpGT       = 1
	jmpEQ       = 2
	jmpGE       = 3
	jmpLT       = 4
	jmpNE       = 5
	jmpLE       = 6
	jmpJMP      = 7
)

var jmpMap = map[string]asm{
	"null": jmpNull,
	"JGT":  jmpGE,
	"JEQ":  jmpEQ,
	"JGE":  jmpGE,
	"JLT":  jmpLT,
	"JNE":  jmpNE,
	"JLE":  jmpLE,
	"JMP":  jmpJMP,
}

////////////////////////////////////////////////////////////////////////////////
// Literals for compute part of instruction.
// Need to be shift six to the left.

const (
	// A = 0
	cmp0       asm = 42
	cmp1           = 63
	cmpNeg1        = 58
	cmpD           = 12
	cmpA           = 48
	cmpNotD        = 13
	cmpNotA        = 49
	cmpNegD        = 15
	cmpNegA        = 51
	cmpDInc        = 31
	cmpAInc        = 55
	cmpDDec        = 14
	cmpADec        = 50
	cmpDPlusA      = 2
	cmpDMinusA     = 19
	cmpAMinusD     = 7
	cmpDAndA       = 0
	cmpDOrA        = 21

	// A = 1
	cmpM       = 112
	cmpNotM    = 113
	cmpNegM    = 79
	cmpMInc    = 119
	cmpMDec    = 114
	cmpDPlusM  = 66
	cmpDMinusM = 83
	cmpMMinusD = 71
	cmpDAndM   = 64
	cmpDOrM    = 85
)

var cmpMap = map[string]asm{
	// A = 0
	"0":   cmp0 << 6,
	"1":   cmp1 << 6,
	"-1":  cmpNeg1 << 6,
	"D":   cmpD << 6,
	"A":   cmpA << 6,
	"!D":  cmpNotD << 6,
	"!A":  cmpNotA << 6,
	"-D":  cmpNegD << 6,
	"-A":  cmpNegA << 6,
	"D+1": cmpDInc << 6,
	"A+1": cmpAInc << 6,
	"D-1": cmpDDec << 6,
	"A-1": cmpADec << 6,
	"D+A": cmpDPlusA << 6,
	"D-A": cmpDMinusA << 6,
	"A-D": cmpAMinusD << 6,
	"D&A": cmpDAndA << 6,
	"D|A": cmpDOrA << 6,

	// A = 1
	"M":   cmpM << 6,
	"!M":  cmpNotM << 6,
	"-M":  cmpNegM << 6,
	"M+1": cmpMInc << 6,
	"M-1": cmpMDec << 6,
	"D+M": cmpDPlusM << 6,
	"D-M": cmpDMinusM << 6,
	"M-D": cmpMMinusD << 6,
	"D&M": cmpDAndM << 6,
	"D|M": cmpDOrM << 6,
}
