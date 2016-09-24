package components

import "os"
import "io/ioutil"

func Assemble(in string, out *os.File) error {
	b, err := ioutil.ReadFile(in)

	if err != nil {
		return err
	}

	input := string(b)

	// Create a Lexer, it will kick off the lexing in a go routine
	lexChan := newLexer(input)

	// Create a parser, and hand it the output from the lexer.  It will
	// run the first pass (building symbol table) and return any errors
	// found by the lexer.  i.e. the lexr and the first pass will be run
	// concurrently, but since we cannot move onto the second pass
	// (parsing the tokens and writing the file) until that's complete,
	// there's no benefit running the second phase concurrently.
	parser, errs := newParser(lexChan)

	if errs != nil {
		return errs.asError()
	}

	// Note that the parser will not write to the file if any errors
	// were found during the lexing phase, and will stop writing after
	// the first error that it encounters.  It is up to the calling
	// routine to tidy the file.
	return parser.run(out).asError()
}
