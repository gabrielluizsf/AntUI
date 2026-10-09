package css

// routeDecl sends one resolved declaration to its domain's handler, one
// case label per property the engine knows. Unknown properties fall to the
// default, which reports them the way a browser reports a misspelling.
func routeDecl(e *declEnv) []string {
	switch e.prop {
	case "background", "background-color", "background-image", "background-position", "background-size", "background-repeat", "background-clip", "background-origin", "background-attachment":
		return applyBackgroundDecl(e)
	case "color", "opacity":
		return applyColorDecl(e)
	case "font-size", "font-weight", "font-style", "font-family", "line-height", "letter-spacing", "word-spacing":
		return applyFontDecl(e)
	case "text-align", "text-transform", "text-decoration", "text-decoration-line", "white-space", "overflow-wrap", "word-wrap", "text-overflow", "vertical-align":
		return applyTextDecl(e)
	case "width", "min-width", "max-width", "height", "min-height", "max-height", "box-sizing":
		return applyBoxDecl(e)
	case "margin", "margin-top", "margin-right", "margin-bottom", "margin-left", "padding", "padding-top", "padding-right", "padding-bottom", "padding-left":
		return applySpacingDecl(e)
	case "border", "border-width", "border-style", "border-color", "border-top", "border-right", "border-bottom", "border-left", "border-radius":
		return applyBorderDecl(e)
	case "display", "position", "top", "right", "bottom", "left", "z-index":
		return applyLayoutDecl(e)
	case "flex-direction", "flex-wrap", "flex-flow":
		return applyFlexDecl(e)
	case "justify-content", "align-items", "align-self", "align-content", "gap", "row-gap", "column-gap":
		return applyAlignDecl(e)
	case "order", "flex", "flex-grow", "flex-shrink", "flex-basis":
		return applyFlexItemDecl(e)
	case "grid-template-columns", "grid-template-rows", "grid-template-areas", "grid-auto-flow", "justify-items", "justify-self":
		return applyGridDecl(e)
	case "grid-column", "grid-row", "grid-column-start", "grid-column-end", "grid-row-start", "grid-row-end", "grid-area":
		return applyGridPlacementDecl(e)
	case "columns", "column-count", "column-width", "column-fill", "column-rule", "column-rule-width", "column-rule-style", "column-rule-color":
		return applyColumnsDecl(e)
	case "break-inside", "page-break-inside", "float", "clear":
		return applyFlowDecl(e)
	case "overflow", "overflow-x", "overflow-y":
		return applyOverflowDecl(e)
	case "visibility", "pointer-events", "cursor":
		return applyVisibilityDecl(e)
	case "box-shadow", "text-shadow", "filter", "backdrop-filter", "transform", "transform-origin":
		return applyEffectDecl(e)
	case "outline", "outline-width", "outline-style", "outline-color", "outline-offset":
		return applyOutlineDecl(e)
	case "transition", "transition-property", "transition-duration", "transition-timing-function", "transition-delay":
		return applyTransitionDecl(e)
	case "animation", "animation-name", "animation-duration", "animation-timing-function", "animation-delay", "animation-iteration-count", "animation-direction", "animation-fill-mode":
		return applyAnimationDecl(e)
	}
	return []string{fmtErrf("ignoring unknown property %q", e.prop).Error()}
}
