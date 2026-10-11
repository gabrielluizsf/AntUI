package template

// The icons the built-in widgets draw themselves, by the name they are
// registered under. A widget that draws one of these asks for it by name, so a
// developer who replaces it changes what every widget of that kind looks like
// at once.
const (
	// IconCheck is the tick inside a checked checkbox.
	IconCheck = "check"
	// IconDot is the filled centre of a selected radio button.
	IconDot = "dot"
	// IconCaret is the arrow that says a dropdown opens downwards.
	IconCaret = "caret"
	// IconCaretUp is the same arrow the other way up, which is what a date
	// picker shows while its calendar is open.
	IconCaretUp = "caret-up"
)
