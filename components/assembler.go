package components

import (
	"io/ioutil"
	"os"
)

// Will take a Hack assembler file (.asm) and writes to out, the
// "binary" machine codes.
func Assemble(in string, out *os.File) error {
	b, err := ioutil.ReadFile(in)

	if err != nil {
		return err
	}

	input := string(b)

	// Create a Lexer, it will kick off the lexing in a go routine
	lexChan := kickOff(input)

	// Create a parser, and hand it the output from the lexer.  It will
	// run the first pass (building symbol table) and return any errors
	// found by the lexer.  i.e. the lexr and the first pass will be run
	// concurrently, but since we cannot move onto the second pass
	// (parsing the tokens and writing the file) until that's complete,
	// there's no benefit to running the second phase concurrently.
	parser, errs := newParser(lexChan)

	if errs != nil {
		return errs.asError()
	}

	// Note that the parser will stop writing after the first error that
	// it encounters.  It is up to the calling routine to tidy the file.
	return parser.run(out).asError()
}
