package css

// specificity returns the selector's [a, b, c] triple: classes and
// pseudo-classes, elements, universals. A higher triple wins the cascade.
func (s Selector) specificity() (a, b, c int) {
	b += len(s.Classes)
	for st := State(1); st <= StateChecked; st <<= 1 {
		if s.State&st != 0 {
			b++
		}
	}
	if !s.All && s.Tag != "" {
		c = 1
	}
	return a, b, c
}
