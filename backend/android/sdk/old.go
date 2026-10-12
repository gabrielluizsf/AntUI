package sdk

func (m *Missing) old(what, have, want string) {
	m.Problems = append(m.Problems, Problem{What: what, Have: have, Want: want})
}
