package css

// applyInitial resets a property to the CSS initial value. Where the initial
// and the engine's "unset" carry the same shape — a zero colour meaning the
// theme's ink, width auto — they stay the same field value apart from colour,
// which keeps its own sentinel so templates can tell "theme ink" apart from a
// declared transparent.
func applyInitial(st *Style, prop string) {
	switch prop {
	case "color", "opacity", "background", "background-color", "background-image",
		"background-position", "background-size", "background-repeat",
		"background-clip", "background-origin", "background-attachment":
		initialSurface(st, prop)
	case "font-size", "text-align", "font-weight", "font-style", "font-family",
		"line-height", "letter-spacing", "word-spacing", "text-transform",
		"text-decoration", "text-decoration-line", "white-space", "overflow-wrap",
		"word-wrap", "text-overflow", "vertical-align":
		initialText(st, prop)
	case "width", "min-width", "height", "min-height", "max-width", "max-height",
		"margin", "margin-top", "margin-right", "margin-bottom", "margin-left",
		"padding", "padding-top", "padding-right", "padding-bottom", "padding-left",
		"box-sizing":
		initialBox(st, prop)
	case "border", "border-width", "border-style", "border-color", "border-radius":
		initialBorder(st, prop)
	case "display", "position", "top", "right", "bottom", "left", "z-index",
		"overflow", "overflow-x", "overflow-y", "visibility", "pointer-events",
		"cursor", "float", "clear", "break-inside", "page-break-inside":
		initialLayout(st, prop)
	case "flex-direction", "flex-wrap", "flex-flow", "justify-content",
		"align-items", "align-self", "align-content", "gap", "row-gap",
		"column-gap", "order", "flex", "flex-grow", "flex-shrink", "flex-basis":
		initialFlex(st, prop)
	case "grid-template-columns", "grid-template-rows", "grid-template-areas",
		"grid-auto-flow", "justify-items", "justify-self", "grid-column",
		"grid-column-start", "grid-column-end", "grid-row", "grid-row-start",
		"grid-row-end", "grid-area":
		initialGrid(st, prop)
	case "columns", "column-count", "column-width", "column-fill", "column-rule",
		"column-rule-width", "column-rule-style", "column-rule-color":
		initialColumns(st, prop)
	case "box-shadow", "text-shadow", "outline", "outline-width", "outline-style",
		"outline-color", "outline-offset", "filter", "backdrop-filter",
		"transform", "transform-origin":
		initialEffects(st, prop)
	case "transition", "transition-property", "transition-duration",
		"transition-timing-function", "transition-delay", "animation",
		"animation-name", "animation-duration", "animation-timing-function",
		"animation-delay", "animation-iteration-count", "animation-direction",
		"animation-fill-mode":
		initialMotion(st, prop)
	}
}
