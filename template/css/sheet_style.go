package css

// Style computes the winning declarations for an element: matching rules are
// folded by importance, specificity and source order, exactly as a web browser
// would. Custom properties (the --foo kind) are resolved first, then normal
// declarations — and every !important declaration applies last, overwriting
// anything from the same cascade scope. baseWidth is the window width, which
// media queries test against and percentage/viewport lengths scale with; the
// window it draws is taken to be as tall as it is wide. [Sheet.StyleViewport]
// hands the cascade a real window instead. The method appends warnings about
// vendor-prefixed properties, unknown properties and unrecognised media
// conditions to the sheet's warning list.
func (sh *Sheet) Style(tag string, classes []string, state State, baseWidth int) Style {
	return sh.StyleViewport(tag, classes, state, Viewport{Width: baseWidth, Height: baseWidth})
}
