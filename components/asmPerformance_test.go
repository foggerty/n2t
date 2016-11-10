// I've added this file purely for performance testing/benchmarking;
// i.e. this is not a test.
package components

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const inputAsm string = "../assembler/Pong.asm"

// BenchmarkPong does benchmarking.
// Oh god linting errors are well meaning pain....
func BenchmarkPong(b *testing.B) {
	if _, err := os.Stat(inputAsm); err != nil {
		msg := fmt.Sprintf("Cannot find input file!\n(%s)\n", inputAsm)
		panic(msg)
	}

	for i := 0; i < b.N; i++ {
		run()
	}
}

func run() {
	// load the input
	in, _ := ioutil.ReadFile(inputAsm)

	// create a lexer/parser
	lexes := StartLexingAsm(string(in))
	parser := NewParser(lexes)

	// read off and dump the output
	for {
		_, ok := <-parser.Output

		if !ok {
			break
		}
	}
}
