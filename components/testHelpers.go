package components

func (act Asmlexine) misMatch(exp Asmlexine) bool {
	return act.Instruction != exp.Instruction || act.Value != exp.Value
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
