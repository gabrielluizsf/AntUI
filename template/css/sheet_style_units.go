package css

// StyleUnits is [Sheet.StyleViewport] with the rest of the measurement
// context: the font sizes for em/rem, and the window the viewport units and
// percentages resolve against. The context's Font is also the cascade hook —
// a font-size declaration living inside the rules updates Font before later
// em lengths resolve.
func (sh *Sheet) StyleUnits(tag string, classes []string, state State, ctx Units) Style {
	return sh.styleUnits(tag, classes, state, ctx, true)
}
