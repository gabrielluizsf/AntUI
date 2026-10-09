package css

// Media is the window a rule lives inside. A stylesheet may list several
// media queries — comma- and or-separated — and the rule applies while any of
// them holds. Every condition the canvas can answer is evaluated against the
// viewport of the frame being drawn: the four min/max edges, orientation off
// the window's two sides, resolution off the display's scale and
// prefers-color-scheme off the color scheme the system paints in. A condition
// with no sensor behind it (a media type such as print) is parsed, warned
// about, and treated as satisfied so the rule is never dropped on a canvas —
// and so is one whose answer the system withheld. Negation flips a query's
// window.
type Media struct {
	Queries []MediaQuery
}
