package template

// The built-in icons, written the way an icon set is: a twenty-four unit square
// of a drawing, stroked rather than filled so the line weight stays even at any
// size, with rounded ends and corners so a twelve-pixel tick does not look
// chipped, and painted in `currentColor` so one drawing serves every colour a
// widget wants to say it in.
//
// The stroke width is why a tick drawn here does not look like the four straight
// lines it replaces: a stroked path is one line of a given weight, so it is
// equally thick along the diagonal and at the corner, which is what a tick is.
var builtinIcons = map[string]string{
	IconCheck: `<svg viewBox="0 0 24 24">
		<path d="M5 12.5 L9.5 17 L19 7" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
	IconDot: `<svg viewBox="0 0 24 24">
		<circle cx="12" cy="12" r="6" fill="currentColor"/>
	</svg>`,
	IconCaret: `<svg viewBox="0 0 24 24">
		<path d="M5 9.5 L12 16 L19 9.5" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
	IconCaretUp: `<svg viewBox="0 0 24 24">
		<path d="M5 14.5 L12 8 L19 14.5" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
	IconCaretLeft: `<svg viewBox="0 0 24 24">
		<path d="M14.5 5 L8 12 L14.5 19" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
	IconCaretRight: `<svg viewBox="0 0 24 24">
		<path d="M9.5 5 L16 12 L9.5 19" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
}
