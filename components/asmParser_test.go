package components

import (
	"fmt"
	"testing"
)

const mismatchLength = "Mismatched results (%s).  Expected %d instructions, got %d."
const mismatchInstruction = "Mismatched instruction (%s).  Expected %s, got %s."

type parseTest struct {
	name         string
	instructions []asmLexeme
	expected     []string
}

var tests = []parseTest{
	{"Single A-Instruction",
		[]asmLexeme{
			{instruction: asmAINSTRUCT, lineNum: 1, value: "7"},
			{instruction: asmEOL, lineNum: 1, value: ""}},
		[]string{"0000000000000111"}},
	{"Single COMP instruction",
		[]asmLexeme{
			{instruction: asmDEST, lineNum: 1, value: "D"},
			{instruction: asmCOMP, lineNum: 1, value: "D+1"},
			{instruction: asmJUMP, lineNum: 1, value: "JEQ"},
			{instruction: asmEOF, lineNum: 1, value: ""}},
		[]string{"1110011111010010"}},
}

func Test(t *testing.T) {
	for _, tst := range tests {
		fmt.Printf("Testing: %s\n", tst.name)

		c := newChannel(tst.instructions)
		p := NewParser(c)
		results := collectResults(p)

		compare(t, tst, results)
	}
}

func collectResults(p AsmParser) []string {
	var results []string

	for {
		asm, ok := <-p.Output

		if !ok {
			break
		}

		results = append(results, asm)
	}

	return results
}

func compare(t *testing.T, pt parseTest, results []string) {
	lenExpected := len(pt.expected)
	lenActual := len(results)

	if lenExpected != lenActual {
		t.Errorf(mismatchLength, pt.name, lenExpected, lenActual)
		return
	}

	for i := range pt.expected {
		expected := pt.expected[i]
		actual := results[i]

		if expected != actual {
			t.Errorf(mismatchInstruction, pt.name, expected, actual)
		}
	}
}

func newChannel(asm []asmLexeme) chan asmLexeme {
	c := make(chan asmLexeme)

	go func() {
		for _, lex := range asm {
			c <- lex
		}

		close(c)
	}()

	return c
}
