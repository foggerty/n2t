package components

import "testing"

// Note to self: defining (an anonymous?) type and creating a literal
// instance of it at the same time.
var asmTest = []struct {
	name     string
	input    string
	expected []Asmlexine
}{
	{"Null input", "",
		[]Asmlexine{{Instruction: asmEOF, Value: ""}}},
	{"Lots o spaces", "        ",
		[]Asmlexine{{Instruction: asmEOF, Value: ""}}},
	{"Tabs 'n things", "    \t  \n\n\t   \t\n   \n\n     \t\t   \t \n \n",
		[]Asmlexine{{Instruction: asmEOF, Value: ""}}},
	{"Label only", "   \n\t   (LOOP)",
		[]Asmlexine{
			{Instruction: asmLABEL, Value: "(LOOP)"},
			{Instruction: asmEOF, Value: ""}}},
	{"A-Instruction only", "@abc123",
		[]Asmlexine{
			{Instruction: asmAINSTRUCT, Value: "@abc123"},
			{Instruction: asmEOF, Value: ""}}},
	{"Single comp instruction", "compy",
		[]Asmlexine{
			{Instruction: asmCOMP, Value: "compy"},
			{Instruction: asmEOF, Value: ""}}},
	{"Single full instruction", "dest=comp;jmp",
		[]Asmlexine{
			{Instruction: asmDEST, Value: "dest"},
			{Instruction: asmCOMP, Value: "comp"},
			{Instruction: asmJUMP, Value: "jmp"},
			{Instruction: asmEOF, Value: ""}}},
}

func TestTheLot(t *testing.T) {
	for _, asmT := range asmTest {
		// new lexer
		_, items := NewLexer(asmT.input)
		results := make([]Asmlexine, 0)

		// collect results
		for res := range items {
			results = append(results, res)
		}

		// run tests
		checkResults(t, asmT.name, asmT.expected, results)
	}
}

func checkResults(t *testing.T, name string, expected []Asmlexine, actual []Asmlexine) {
	const lengthMismatch = "%s: wa s expecting to get %d tokens, but got %d."
	const mismatchedToken = "%s: mismatch, got %q but expected %q."

	lenExpected := len(expected)
	lenActual := len(actual)

	if lenExpected != lenActual {
		t.Errorf(lengthMismatch, name, lenExpected, lenActual)
	}

	for i := 0; i < min(lenExpected, lenActual); i++ {
		diff := expected[i].misMatch(actual[i])

		if diff {
			t.Errorf(mismatchedToken, name, actual[i], expected[i])
		}
	}
}
