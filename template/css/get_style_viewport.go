package css

// GetStyleViewport computes the winning style for one widget against a whole
// window: body for the window itself, or one of the Role tags. Classes are
// the widget's own, combined with the class table's entries for its kind.
// Media queries test both edges of the viewport along with its density and
// the system's color scheme, and state carries the pseudo-classes the widget
// is in. The result is cached — but only while the window stays the same,
// because media queries are read off the window, so a resize or a theme the
// system changed makes every cached style stale and the table starts over on
// the next lookup. A caller that paints every widget each frame therefore
// computes each style once per window, not once ever.
func (c *CSSClasses) GetStyleViewport(q StyleQuery) Style {
	if c.sheet == nil {
		return Style{}
	}
	if c.vp != q.Viewport {
		clear(c.specs)
		c.vp = q.Viewport
	}
	// The cache is keyed by what decides the answer: the widget kind, its
	// state, and the window its media queries and lengths resolve against. A
	// caller that adds classes of its own is asking for one widget to differ
	// from the table's, and is looked up rather than cached — a map cannot key
	// a list without building a string for it every frame, which is the very
	// thing the cache is here to avoid. The template, which is the caller that
	// matters, asks for the table's own classes and hits the cache every time.
	if len(q.Classes) == 0 {
		key := styleKey{tag: q.Tag, state: q.State, width: q.Viewport.Width, height: q.Viewport.Height}
		if st, ok := c.specs[key]; ok {
			return st
		}
		st := c.sheet.StyleViewport(q.Tag, c.classesOf(q.Tag), q.State, q.Viewport)
		if c.specs != nil {
			c.specs[key] = st
		}
		return st
	}
	all := append(append([]string(nil), c.classesOf(q.Tag)...), q.Classes...)
	return c.sheet.StyleViewport(q.Tag, all, q.State, q.Viewport)
}
