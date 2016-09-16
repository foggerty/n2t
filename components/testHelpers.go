package components

func (act AsmLexeme) misMatch(exp AsmLexeme) bool {
	return act.instruction != exp.instruction || act.value != exp.value
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
