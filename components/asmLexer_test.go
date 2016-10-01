package components

import "testing"
import "fmt"

// Note to self: defining (an anonymous?) type and creating a literal
// instance of it at the same time.
var asmTest = []struct {
	name     string
	input    string
	expected []asmLexeme
}{
	{"Null input", "",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOF, value: ""}}},

	{"Lots o spaces", "        ",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOF, value: ""}}},

	{"Lots o comments 1", "// La la\n//woo woo    \n     //Wheee!",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmEOL, value: ""},
			{lineNum: 3, instruction: asmEOF, value: ""}}},

	{"Two A-Instructions", "@1\n@2",
		[]asmLexeme{
			{lineNum: 1, instruction: asmAINSTRUCT, value: "1"},
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmAINSTRUCT, value: "2"},
			{lineNum: 2, instruction: asmEOF, value: ""}}},

	{"Tabs 'n things", "    \t  \n\n\t   \t\n   \n\n     \t\t   \t \n \n",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmEOL, value: ""},
			{lineNum: 3, instruction: asmEOL, value: ""},
			{lineNum: 4, instruction: asmEOL, value: ""},
			{lineNum: 5, instruction: asmEOL, value: ""},
			{lineNum: 6, instruction: asmEOL, value: ""},
			{lineNum: 7, instruction: asmEOL, value: ""},
			{lineNum: 8, instruction: asmEOF, value: ""}}},

	{"Label only", "\n\t(LOOP)",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmLABEL, value: "LOOP"},
			{lineNum: 2, instruction: asmEOF, value: ""}}},

	{"A-Instruction only", "@abc123",
		[]asmLexeme{
			{lineNum: 1, instruction: asmAINSTRUCT, value: "abc123"},
			{lineNum: 1, instruction: asmEOF, value: ""}}},

	{"Single comp instruction", "\n\n  \t D+1   \r\n",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmEOL, value: ""},
			{lineNum: 3, instruction: asmERROR, value: "  \t D+1   "},
			{lineNum: 3, instruction: asmEOL, value: ""},
			{lineNum: 4, instruction: asmEOF, value: ""}}},

	{"Single full instruction", "AMD=D+1;JMP",
		[]asmLexeme{
			{lineNum: 1, instruction: asmDEST, value: "AMD"},
			{lineNum: 1, instruction: asmCOMP, value: "D+1"},
			{lineNum: 1, instruction: asmJUMP, value: "JMP"},
			{lineNum: 1, instruction: asmEOF, value: ""}}},

	{"Single full instruction with newlines and comments", "//moose \n   //wibble\n\nD=D&M;JLT\n\n",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmEOL, value: ""},
			{lineNum: 3, instruction: asmEOL, value: ""},
			{lineNum: 4, instruction: asmDEST, value: "D"},
			{lineNum: 4, instruction: asmCOMP, value: "D&M"},
			{lineNum: 4, instruction: asmJUMP, value: "JLT"},
			{lineNum: 4, instruction: asmEOL, value: ""},
			{lineNum: 5, instruction: asmEOL, value: ""},
			{lineNum: 6, instruction: asmEOF, value: ""}}},

	{"Three instructions.", " D=D+M  \n @123\n  (LOOP)",
		[]asmLexeme{
			{lineNum: 1, instruction: asmDEST, value: "D"},
			{lineNum: 1, instruction: asmCOMP, value: "D+M"},
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmAINSTRUCT, value: "123"},
			{lineNum: 2, instruction: asmEOL, value: ""},
			{lineNum: 3, instruction: asmLABEL, value: "LOOP"},
			{lineNum: 3, instruction: asmEOF, value: ""}}},

	{"Handles Windows EOL", "   \n   \r\n  (LOOP)  \r\n",
		[]asmLexeme{
			{lineNum: 1, instruction: asmEOL, value: ""},
			{lineNum: 2, instruction: asmEOL, value: ""},
			{lineNum: 3, instruction: asmLABEL, value: "LOOP"},
			{lineNum: 3, instruction: asmEOL, value: ""},
			{lineNum: 4, instruction: asmEOF, value: ""},
		}},

	{"Two instructions on one line", "D=D+1    @fred",
		[]asmLexeme{
			{lineNum: 1, instruction: asmDEST, value: "D"},
			{lineNum: 1, instruction: asmCOMP, value: "D+1"},
			{lineNum: 1, instruction: asmERROR, value: "D=D+1    @fred"},
			{lineNum: 1, instruction: asmEOF, value: ""},
		}},
}

func TestTheLot(t *testing.T) {
	for _, asmT := range asmTest {
		fmt.Println("Running test: " + asmT.name)

		// new lexer
		items := kickOff(asmT.input)
		results := make([]asmLexeme, 0)

		// collect results
		for res := range items {
			results = append(results, res)
		}

		// run tests
		checkResults(t, asmT.name, asmT.expected, results)
	}
}

func checkResults(t *testing.T, name string, expected []asmLexeme, actual []asmLexeme) {
	const lengthMismatch = "%s:\nwas expecting to get %d tokens, but got %d."
	const mismatchedToken = "%s:\nExpected %q but got %q."
	const incorrectLineNum = "%s:\nFor instruction %q, expected line number of %d, got %d."

	lenExpected := len(expected)
	lenActual := len(actual)

	if lenExpected != lenActual {
		t.Errorf(lengthMismatch, name, lenExpected, lenActual)
	}

	for i := 0; i < min(lenExpected, lenActual); i++ {
		diff := expected[i].misMatch(actual[i])

		if diff {
			t.Errorf(mismatchedToken, name, expected[i], actual[i])
		}

		eL := expected[i].lineNum
		aL := actual[i].lineNum

		if eL != aL {
			t.Errorf(incorrectLineNum, name, expected[i], eL, aL)
		}
	}
}
