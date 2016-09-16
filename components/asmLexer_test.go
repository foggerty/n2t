package components

import "testing"

// Note to self: defining (an anonymous?) type and creating a literal
// instance of it at the same time.
var asmTest = []struct {
	name     string
	input    string
	expected []AsmLexeme
}{
	{"Null input", "",
		[]AsmLexeme{{instruction: asmEOF, value: ""}}},

	{"Lots o spaces", "        ",
		[]AsmLexeme{{instruction: asmEOF, value: ""}}},

	{"Lots o comments 1", "// La la\n//woo woo    \n     //Wheee!",
		[]AsmLexeme{
			{instruction: asmEOF, value: ""},
		}},

	{"Blank lines", "@1\n@2",
		[]AsmLexeme{
			{instruction: asmAINSTRUCT, value: "1"},
			{instruction: asmEOL, value: ""},
			{instruction: asmAINSTRUCT, value: "2"},
			{instruction: asmEOF, value: ""},
		}},

	{"Tabs 'n things", "    \t  \n\n\t   \t\n   \n\n     \t\t   \t \n \n",
		[]AsmLexeme{
			{instruction: asmEOF, value: ""}}},

	{"Label only", "\n\t(LOOP)",
		[]AsmLexeme{
			{instruction: asmLABEL, value: "LOOP"},
			{instruction: asmEOF, value: ""}}},

	{"A-Instruction only", "@abc123",
		[]AsmLexeme{
			{instruction: asmAINSTRUCT, value: "abc123"},
			{instruction: asmEOF, value: ""}}},

	{"Single comp instruction", "\n\nD+1",
		[]AsmLexeme{
			{instruction: asmERROR, value: "Unknown error at line 3 (D+1)"},
			{instruction: asmEOF, value: ""}}},

	{"Single full instruction", "AMD=D+1;JMP",
		[]AsmLexeme{
			{instruction: asmDEST, value: "AMD"},
			{instruction: asmCOMP, value: "D+1"},
			{instruction: asmJUMP, value: "JMP"},
			{instruction: asmEOF, value: ""}}},

	{"Single full instruction with newlines and comments", "//moose \n   //wibble\n\nD=D&M;JLT\n\n",
		[]AsmLexeme{
			{instruction: asmDEST, value: "D"},
			{instruction: asmCOMP, value: "D&M"},
			{instruction: asmJUMP, value: "JLT"},
			{instruction: asmEOL, value: ""},
			{instruction: asmEOF, value: ""}}},

	{"Three instructions.", " D=D+M  \n @123\n  (LOOP)",
		[]AsmLexeme{
			{instruction: asmDEST, value: "D"},
			{instruction: asmCOMP, value: "D+M"},
			{instruction: asmEOL, value: ""},
			{instruction: asmAINSTRUCT, value: "123"},
			{instruction: asmEOL, value: ""},
			{instruction: asmLABEL, value: "LOOP"},
			{instruction: asmEOF, value: ""}}},
}

func TestTheLot(t *testing.T) {
	for _, asmT := range asmTest {
		// new lexer
		items := newLexer(asmT.input)
		results := make([]AsmLexeme, 0)

		// collect results
		for res := range items {
			results = append(results, res)
		}

		// run tests
		checkResults(t, asmT.name, asmT.expected, results)
	}
}

func checkResults(t *testing.T, name string, expected []AsmLexeme, actual []AsmLexeme) {
	const lengthMismatch = "%s:\n was expecting to get %d tokens, but got %d."
	const mismatchedToken = "%s:\n Expected %q but got %q."

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
	}
}
