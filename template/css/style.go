package css

// Style is the computed drawing style of one widget after the cascade: which
// pixels to fill, how thick its border is, how big its text. Zero values fall
// back to the template's theme — a style does not know the theme, it knows
// only what the stylesheet said, and the Set map records whether it spoke.
// The field groups live in their own types (boxStyle, surfaceStyle, …) so
// each stays a readable size; they are embedded, so every field reads and
// writes as if it were declared right here.
type Style struct {
	boxStyle
	surfaceStyle
	textStyle
	layoutStyle
	flexStyle
	gridStyle
	columnStyle
	effectStyle
	motionStyle

	// Set records which canonical properties the stylesheet mentioned, so a
	// template can decide theme fallbacks.
	Set map[string]bool

	// Custom holds the cascaded custom properties (--foo: …) of the element,
	// so a template can read what the stylesheet named them. var(--foo)
	// references in other declarations are resolved against this map before
	// a value is parsed.
	Custom map[string]string

	// inherit records the properties whose winning cascade value was the
	// inherit keyword (or unset on an inherited property). StyleUnits folds
	// them in from the body's computed style after the cascade; the map is
	// internal to that pass and never meant for the template.
	inherit map[string]bool
}

// Has reports whether a canonical property was set in the cascade.
func (s *Style) Has(prop string) bool { return s.Set[prop] }
