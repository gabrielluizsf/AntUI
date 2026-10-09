package css

// Rule is one selector-list block, optionally constrained by a media query.
type Rule struct {
	Media     Media
	Selectors []Selector
	Decls     []Declaration
	Order     int
}
