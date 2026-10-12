package svg

// applyAttrs reads every attribute an element wrote into s, in whatever order
// the map hands them over: each case touches fields of its own, so the order
// they are read in is not one to depend on. The dash pattern is the exception
// and is left to [Style.with], which reads it after all of them.
func (s *Style) applyAttrs(attrs map[string]string, warn func(string, ...any)) {
	s.applyPaintAttrs(attrs, warn)
	s.applyStrokeAttrs(attrs, warn)
	s.applyTextAttrs(attrs, warn)
}
