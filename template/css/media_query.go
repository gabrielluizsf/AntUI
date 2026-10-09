package css

// MediaQuery is one alternative of an @media list: a window over the
// viewport, with an optional negation. A query with no window matches any
// viewport; that is also how conditions the engine cannot evaluate fall out.
type MediaQuery struct {
	MinWidth     int
	MaxWidth     int
	MinHeight    int
	MaxHeight    int
	HasMin       bool
	HasMax       bool
	HasMinHeight bool
	HasMaxHeight bool

	// Resolution is the window's device density in dots per inch, read
	// from the display's scale. An exact resolution writes both bounds.
	MinDpi    float64
	MaxDpi    float64
	HasMinDpi bool
	HasMaxDpi bool

	// Orientation is the window's shape: Portrait at least as tall as it
	// is wide, landscape the other way round.
	HasOrientation bool
	Portrait       bool

	// Scheme is prefers-color-scheme, which a system that did not say
	// leaves out of the answer entirely.
	HasScheme bool
	Dark      bool

	Negated bool

	// never marks two nested queries that contradicted each other — a
	// portrait window inside a landscape one, a light scheme inside a dark
	// one — and no viewport is on both sides of that.
	never bool
}
