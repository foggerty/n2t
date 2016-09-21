package components

func (act asmLexeme) misMatch(exp asmLexeme) bool {
	return act.instruction != exp.instruction || act.value != exp.value
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
