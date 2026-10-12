package sdk

func (m *Missing) add(what, install string) {
	m.Problems = append(m.Problems, Problem{What: what, Install: install})
}
