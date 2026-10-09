package css

func (t *mathTokens) empty() bool {
	return t.peek() == ""
}
