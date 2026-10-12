package svg

// childRuns reads a `<tspan>`, the way any other element is read, over the
// style the writing around it has: a size or a colour on it is the size or
// colour of what it holds and nothing else.
func (img *Image) childRuns(e *element, st Style, warn func(string, ...any)) []TextRun {
	kid := st.with(e, warn)
	img.resolvePaints(&kid, warn)
	return img.textRuns(e, kid, warn)
}
