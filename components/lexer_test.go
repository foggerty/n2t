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
		[]AsmLexeme{{Instruction: asmEOF, Value: ""}}},

	{"Lots o spaces", "        ",
		[]AsmLexeme{{Instruction: asmEOF, Value: ""}}},

	{"Lots o comments 1", "// La la\n//woo woo    \n     //Wheee!",
		[]AsmLexeme{
			{Instruction: asmEOF, Value: ""},
		}},

	{"Blank lines", "@1\n@2",
		[]AsmLexeme{
			{Instruction: asmAINSTRUCT, Value: "1"},
			{Instruction: asmEOL, Value: ""},
			{Instruction: asmAINSTRUCT, Value: "2"},
			{Instruction: asmEOF, Value: ""},
		}},

	{"Tabs 'n things", "    \t  \n\n\t   \t\n   \n\n     \t\t   \t \n \n",
		[]AsmLexeme{
			{Instruction: asmEOF, Value: ""}}},

	{"Label only", "\n\t(LOOP)",
		[]AsmLexeme{
			{Instruction: asmLABEL, Value: "LOOP"},
			{Instruction: asmEOF, Value: ""}}},

	{"A-Instruction only", "@abc123",
		[]AsmLexeme{
			{Instruction: asmAINSTRUCT, Value: "abc123"},
			{Instruction: asmEOF, Value: ""}}},

	{"Single comp instruction", "\n\nD+1",
		[]AsmLexeme{
			{Instruction: asmERROR, Value: "Unknown error at line 3 (D+1)"},
			{Instruction: asmEOF, Value: ""}}},

	{"Single full instruction", "AMD=D+1;JMP",
		[]AsmLexeme{
			{Instruction: asmDEST, Value: "AMD"},
			{Instruction: asmCOMP, Value: "D+1"},
			{Instruction: asmJUMP, Value: "JMP"},
			{Instruction: asmEOF, Value: ""}}},

	{"Single full instruction with newlines and comments", "//moose \n   //wibble\n\nD=D&M;JLT\n\n",
		[]AsmLexeme{
			{Instruction: asmDEST, Value: "D"},
			{Instruction: asmCOMP, Value: "D&M"},
			{Instruction: asmJUMP, Value: "JLT"},
			{Instruction: asmEOL, Value: ""},
			{Instruction: asmEOF, Value: ""}}},

	{"Three instructions.", " D=D+M  \n @123\n  (LOOP)",
		[]AsmLexeme{
			{Instruction: asmDEST, Value: "D"},
			{Instruction: asmCOMP, Value: "D+M"},
			{Instruction: asmEOL, Value: ""},
			{Instruction: asmAINSTRUCT, Value: "123"},
			{Instruction: asmEOL, Value: ""},
			{Instruction: asmLABEL, Value: "LOOP"},
			{Instruction: asmEOF, Value: ""}}},
}

func TestTheLot(t *testing.T) {
	for _, asmT := range asmTest {
		// new lexer
		_, items := NewLexer(asmT.input)
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
