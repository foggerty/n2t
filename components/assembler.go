/*
  Assembler

  Exists to link the Lexer and Parser together.

  Lexer (asmLexeme) -> Parser (asm) -> Assembler (output file)
*/

package components

import (
	"fmt"
	"os"
)

func Assemble(in string, out *os.File) {
	input := "" // load from disk

	// Create a Lexer
	lexChan := newLexer(input)

	// Create a parser
	asmChan := newParser(lexChan)

	// Write to file
	for i := range asmChan {
		fmt.Fprint(out, i)
	}
}
