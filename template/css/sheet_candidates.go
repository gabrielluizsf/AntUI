package css

// candidate is one matching rule folded into the cascade: its specificity
// triple (classes, elements), its order in the stylesheet, and its
// declarations.
type candidate struct {
	spec  [3]int
	order int
	decls []Declaration
}
