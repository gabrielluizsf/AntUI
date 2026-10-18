package svg

// relativeSelector is one selector of a `:has(…)`: the combinator that joins
// it to the element the `:has` stands on, and the complex selector it opens.
// Written `:has(> rect)` the lead is a child combinator; written `:has(rect)`
// it is a descendant, which is what a relative selector with nothing in front
// means.
type relativeSelector struct {
	lead combinator
	cs   complexSelector
}

// nthExpr is the `An+B` of an `:nth-…`, with the `of S` list a couple of them
// may carry. A sibling stands where `A*k+B` equals its position for some k of
// zero or more.
type nthExpr struct {
	a  int
	b  int
	of []complexSelector
}
