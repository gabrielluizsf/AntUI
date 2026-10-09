package css

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

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

// applyDecl folds one declaration into the style, expanding shorthands,
// resolving var() references and evaluating length formulas (calc/min/max/
// clamp) under the measurement context. Properties the engine does not know
// are ignored and warned about, exactly as a browser ignores and cools its
// heels on them. The context's Font is updated by a font-size declaration so
// later em lengths see it. applyDecl returns the warnings it produced.
func applyDecl(st *Style, set map[string]bool, customs, cascaded map[string]string, d Declaration, ctx *Units) []string {
	prop := strings.ToLower(strings.TrimSpace(d.Prop))
	raw := strings.TrimSpace(d.Raw)
	if prop == "" || raw == "" {
		return nil
	}
	if set == nil {
		st.Set = make(map[string]bool, 16)
		set = st.Set
	}
	note := func() { delete(st.inherit, prop); set[prop] = true }

	// Custom property detection runs before the vendor-prefix rule, because
	// --foo starts with a dash too. Their values are captured for var() to
	// resolve later.
	if strings.HasPrefix(prop, "--") {
		st.Custom[prop] = raw
		return nil
	}
	if strings.HasPrefix(prop, "-") {
		return []string{fmtErrf("ignoring vendor-prefixed property %q", prop).Error()}
	}
	raw = resolveVars(raw, cascaded)
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return nil
	}
	if raw == "initial" || raw == "inherit" || raw == "unset" || raw == "revert" {
		return applyKeyword(st, note, raw, prop)
	}

	switch prop {
	case "background":
		if applyBackground(st, raw, *ctx) {
			note()
			for _, kp := range backgroundKeys {
				set[kp] = true
			}
		}
	case "background-color":
		if c, err := ParseColor(raw); err == nil {
			st.Background = c
			note()
		}
	case "background-image":
		if v, ok := parseBackgroundImage(raw, *ctx); ok {
			st.BackgroundImages = v
			note()
		}
	case "background-position":
		if v, ok := parseBackgroundPosition(raw, *ctx); ok {
			st.BackgroundPos = v
			note()
		}
	case "background-size":
		if v, ok := parseBackgroundSize(raw, *ctx); ok {
			st.BackgroundSize = v
			note()
		}
	case "background-repeat":
		if v, ok := parseBackgroundRepeat(raw); ok {
			st.BackgroundRepeat = v
			note()
		}
	case "background-clip":
		if v, ok := parseBackgroundBox(raw); ok {
			st.BackgroundClip = v
			note()
		}
	case "background-origin":
		if v, ok := parseBackgroundBox(raw); ok {
			st.BackgroundOrigin = v
			note()
		}
	case "background-attachment":
		if v, ok := parseBackgroundAttachment(raw); ok {
			st.BackgroundAttach = v
			note()
		}
	case "color":
		if c, err := ParseColor(raw); err == nil {
			st.Color = c
			note()
		}
	case "opacity":
		if v, ok := parseOpacity(raw); ok {
			st.Opacity = v
			note()
		}
	case "font-size":
		if l, err := parseLengthAt(raw, *ctx); err == nil && !l.IsPct() {
			st.FontSize = l.Resolve(*ctx)
			if st.FontSize <= 0 {
				st.FontSize = DefaultFontSize
			}
			ctx.Font = st.FontSize
			note()
		}
	case "text-align":
		switch raw {
		case "left":
			st.TextAlign = 0
			note()
		case "center":
			st.TextAlign = 1
			note()
		case "right":
			st.TextAlign = 2
			note()
		case "justify":
			st.TextAlign = 3
			note()
		}
	case "font-weight":
		if w, ok := parseFontWeight(raw); ok {
			st.FontWeight = w
			note()
		}
	case "font-style":
		if v, ok := parseFontStyle(raw); ok {
			st.FontStyle = v
			note()
		}
	case "font-family":
		if v := normalizeFamily(raw); v != "" {
			st.FontFamily = v
			note()
		}
	case "line-height":
		if v, ok := parseLineHeight(raw, *ctx); ok {
			st.LineHeight = v
			note()
		}
	case "letter-spacing":
		if raw == "normal" {
			st.LetterSpacing = Zero()
			note()
		} else if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.LetterSpacing = l
			note()
		}
	case "word-spacing":
		if raw == "normal" {
			st.WordSpacing = Zero()
			note()
		} else if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.WordSpacing = l
			note()
		}
	case "text-transform":
		if v, ok := parseTextTransform(raw); ok {
			st.TextTransform = v
			note()
		}
	case "text-decoration", "text-decoration-line":
		if bits, ok := parseTextDecoration(raw); ok {
			st.TextDecoration = bits
			note()
		}
	case "white-space":
		if v, ok := parseWhiteSpace(raw); ok {
			st.WhiteSpace = v
			note()
		}
	case "overflow-wrap", "word-wrap":
		if v, ok := parseOverflowWrap(raw); ok {
			st.OverflowWrap = v
			note()
			set["overflow-wrap"] = true
		}
	case "text-overflow":
		if v, ok := parseTextOverflow(raw); ok {
			st.TextOverflow = v
			note()
		}
	case "vertical-align":
		if v, l, ok := parseVerticalAlign(raw, *ctx); ok {
			st.VerticalAlign = v
			st.BaselineShift = l
			note()
		}
	case "width":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Width = l
			note()
		}
	case "height":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Height = l
			note()
		}
	case "min-width":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MinWidth = l
			note()
		}
	case "max-width":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MaxWidth = l
			note()
		}
	case "min-height":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MinHeight = l
			note()
		}
	case "max-height":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.MaxHeight = l
			note()
		}
	case "margin":
		if sides, err := parseFourAt(raw, *ctx); err == nil {
			st.Margin = sides
			note()
		}
	case "margin-top":
		setLenSide(st.Margin[:], "margin", set, 0, raw, *ctx)
	case "margin-right":
		setLenSide(st.Margin[:], "margin", set, 1, raw, *ctx)
	case "margin-bottom":
		setLenSide(st.Margin[:], "margin", set, 2, raw, *ctx)
	case "margin-left":
		setLenSide(st.Margin[:], "margin", set, 3, raw, *ctx)
	case "padding":
		if sides, err := parseFourAt(raw, *ctx); err == nil {
			st.Padding = sides
			note()
		}
	case "padding-top":
		setLenSide(st.Padding[:], "padding", set, 0, raw, *ctx)
	case "padding-right":
		setLenSide(st.Padding[:], "padding", set, 1, raw, *ctx)
	case "padding-bottom":
		setLenSide(st.Padding[:], "padding", set, 2, raw, *ctx)
	case "padding-left":
		setLenSide(st.Padding[:], "padding", set, 3, raw, *ctx)
	case "box-sizing":
		switch raw {
		case "border-box":
			st.BoxSizing = true
		case "content-box":
			st.BoxSizing = false
		}
		note()
	case "border":
		applyBorder(st, set, [4]int{0, 1, 2, 3}, raw, note)
	case "border-width":
		if v, ok := parseBorderWidths(raw); ok {
			st.BorderWidth = v
			note()
		}
	case "border-style":
		if v, ok := parseBorderStyles(raw); ok {
			st.BoxStyle = v
			note()
		}
	case "border-color":
		if v, ok := parseBorderColors(raw); ok {
			st.BoxColor = v
			note()
		}
	case "border-top":
		applyBorder(st, set, [4]int{0, -1, -1, -1}, raw, note)
	case "border-right":
		applyBorder(st, set, [4]int{-1, 1, -1, -1}, raw, note)
	case "border-bottom":
		applyBorder(st, set, [4]int{-1, -1, 2, -1}, raw, note)
	case "border-left":
		applyBorder(st, set, [4]int{-1, -1, -1, 3}, raw, note)
	case "border-radius":
		if rx, ry, ok := parseRadiiXY(raw); ok {
			st.Radius = rx
			st.RadiusY = ry
			note()
		}
	case "display":
		if v, ok := parseDisplay(raw); ok {
			st.Display = v
			note()
		}
	case "position":
		if v, ok := parsePosition(raw); ok {
			st.Position = v
			note()
		}
	case "top":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Top = l
			note()
		}
	case "right":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Right = l
			note()
		}
	case "bottom":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Bottom = l
			note()
		}
	case "left":
		if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.Left = l
			note()
		}
	case "z-index":
		if v, ok := parseZIndex(raw); ok {
			st.ZIndex = v
			note()
		}
	case "flex-direction":
		if v, ok := parseFlexDirection(raw); ok {
			st.FlexDirection = v
			note()
		}
	case "flex-wrap":
		if v, ok := parseFlexWrap(raw); ok {
			st.FlexWrap = v
			note()
		}
	case "flex-flow":
		words := splitWords(raw)
		if len(words) == 0 || len(words) > 2 {
			break
		}
		ok := true
		for _, w := range words {
			if v, good := parseFlexDirection(w); good {
				st.FlexDirection = v
				continue
			}
			if v, good := parseFlexWrap(w); good {
				st.FlexWrap = v
				continue
			}
			ok = false
		}
		if ok {
			note()
		}
	case "justify-content":
		if v, ok := parseJustifyContent(raw); ok {
			st.JustifyContent = v
			if grid, good := parseGridJustifyContent(raw); good {
				st.GridJustifyContent = grid
			}
			note()
		}
	case "align-items":
		if v, ok := parseAlignItems(raw); ok {
			st.AlignItems = v
			note()
		}
	case "align-self":
		if v, ok := parseAlignSelf(raw); ok {
			st.AlignSelf = v
			note()
		}
	case "align-content":
		if v, ok := parseAlignContent(raw); ok {
			st.AlignContent = v
			note()
		}
	case "gap":
		if r, c, ok := parseGap(raw, *ctx); ok {
			st.RowGap, st.ColumnGap = r, c
			note()
			set["row-gap"] = true
			set["column-gap"] = true
		}
	case "row-gap":
		if r, _, ok := parseGap(raw, *ctx); ok {
			st.RowGap = r
			note()
			set["gap"] = true
		}
	case "column-gap":
		if _, c, ok := parseGap(raw, *ctx); ok {
			st.ColumnGap = c
			note()
			set["gap"] = true
		}
	case "order":
		if v, ok := parseOrder(raw); ok {
			st.Order = v
			note()
		}
	case "flex":
		if g, s, b, ok := parseFlex(raw, *ctx); ok {
			st.FlexGrow, st.FlexShrink, st.FlexBasis = g, s, b
			note()
			set["flex-grow"] = true
			set["flex-shrink"] = true
			set["flex-basis"] = true
		}
	case "flex-grow":
		if v, ok := parseFlexNumber(raw); ok {
			st.FlexGrow = v
			note()
		}
	case "flex-shrink":
		if v, ok := parseFlexNumber(raw); ok {
			st.FlexShrink = v
			note()
		}
	case "flex-basis":
		if v, ok := parseFlexBasis(raw, *ctx); ok {
			st.FlexBasis = v
			note()
		}
	case "grid-template-columns":
		if v, ok := parseGridTemplate(raw, *ctx); ok {
			st.GridTemplateColumns = v
			note()
		}
	case "grid-template-rows":
		if v, ok := parseGridTemplate(raw, *ctx); ok {
			st.GridTemplateRows = v
			note()
		}
	case "grid-template-areas":
		if v, ok := parseGridTemplateAreas(raw); ok {
			st.GridTemplateAreas = v
			note()
		}
	case "grid-auto-flow":
		if flow, dense, ok := parseGridAutoFlow(raw); ok {
			st.GridAutoFlow = flow
			st.GridAutoFlowDense = dense
			note()
		}
	case "justify-items":
		if v, ok := parseJustifyItems(raw); ok {
			st.JustifyItems = v
			note()
		}
	case "justify-self":
		if v, ok := parseJustifySelf(raw); ok {
			st.JustifySelf = v
			note()
		}
	case "grid-column":
		if v, ok := parseGridPlacement(raw); ok {
			st.GridColumn = v
			note()
			set["grid-column-start"] = true
			set["grid-column-end"] = true
		}
	case "grid-row":
		if v, ok := parseGridPlacement(raw); ok {
			st.GridRow = v
			note()
			set["grid-row-start"] = true
			set["grid-row-end"] = true
		}
	case "grid-column-start", "grid-column-end":
		if v, ok := parseGridLine(raw); ok {
			if prop == "grid-column-start" {
				st.GridColumn.Start = v
			} else {
				st.GridColumn.End = v
			}
			note()
			set["grid-column"] = true
		}
	case "grid-row-start", "grid-row-end":
		if v, ok := parseGridLine(raw); ok {
			if prop == "grid-row-start" {
				st.GridRow.Start = v
			} else {
				st.GridRow.End = v
			}
			note()
			set["grid-row"] = true
		}
	case "grid-area":
		if column, row, ok := parseGridArea(raw); ok {
			st.GridColumn = column
			st.GridRow = row
			note()
			set["grid-column"] = true
			set["grid-row"] = true
			set["grid-column-start"] = true
			set["grid-column-end"] = true
			set["grid-row-start"] = true
			set["grid-row-end"] = true
		}
	case "columns":
		if n, l, ok := parseColumns(raw, *ctx); ok {
			st.ColumnCount, st.ColumnWidth = n, l
			note()
			set["column-count"] = true
			set["column-width"] = true
		}
	case "column-count":
		if v, ok := parseColumnCount(raw); ok {
			st.ColumnCount = v
			note()
		}
	case "column-width":
		if v, ok := parseColumnWidth(raw, *ctx); ok {
			st.ColumnWidth = v
			note()
		}
	case "column-fill":
		if v, ok := parseColumnFill(raw); ok {
			st.ColumnFill = v
			note()
		}
	case "column-rule":
		applyColumnRule(st, set, raw, note)
	case "column-rule-width":
		if l, err := parseLength(raw); err == nil && l.u == unitPx {
			st.ColumnRuleWidth = max(l.Px(0), 0)
			note()
		}
	case "column-rule-style":
		if v, ok := parseBorderStyle(raw); ok {
			st.ColumnRuleStyle = v
			note()
		}
	case "column-rule-color":
		if c, err := ParseColor(raw); err == nil {
			st.ColumnRuleColor = c
			note()
		}
	case "break-inside", "page-break-inside":
		if v, ok := parseBreakInside(raw); ok {
			st.BreakInside = v
			note()
			set["break-inside"] = true
		}
	case "float":
		if v, ok := parseFloat(raw); ok {
			st.Float = v
			note()
			if v != FloatNone {
				return []string{fmtErrf("%s is read but not applied: this template's flow has no line box to float inside", prop).Error()}
			}
		}
	case "clear":
		if v, ok := parseClear(raw); ok {
			st.Clear = v
			note()
			if v != ClearNone {
				return []string{fmtErrf("%s is read but not applied: nothing here floats, so there is nothing to be pushed past", prop).Error()}
			}
		}
	case "overflow":
		if v, ok := parseOverflow(raw); ok {
			st.Overflow = [2]uint8{v, v}
			note()
			set["overflow-x"] = true
			set["overflow-y"] = true
		}
	case "overflow-x":
		if v, ok := parseOverflow(raw); ok {
			st.Overflow[0] = v
			note()
			set["overflow"] = true
		}
	case "overflow-y":
		if v, ok := parseOverflow(raw); ok {
			st.Overflow[1] = v
			note()
			set["overflow"] = true
		}
	case "visibility":
		if v, ok := parseVisibility(raw); ok {
			st.Visibility = v
			note()
		}
	case "pointer-events":
		if v, ok := parsePointerEvents(raw); ok {
			st.PointerEvents = v
			note()
		}
	case "cursor":
		if v, ok := parseCursor(raw); ok {
			st.Cursor = v
			note()
		}
	case "box-shadow":
		if sh, ok := parseShadows(raw, *ctx, true, 4); ok {
			st.BoxShadow = sh
			note()
		}
	case "text-shadow":
		if sh, ok := parseShadows(raw, *ctx, false, 3); ok {
			st.TextShadow = sh
			note()
		}
	case "outline":
		if applyOutline(st, raw, *ctx) {
			note()
		}
	case "outline-width":
		if l, err := parseLengthAt(raw, *ctx); err == nil && !l.IsPct() && !l.Auto() && !l.None() {
			st.OutlineWidth = l.Resolve(*ctx)
			note()
		}
	case "outline-style":
		if s, ok := outlineStyle(raw); ok {
			st.OutlineStyle = s
			note()
		}
	case "outline-color":
		if c, err := ParseColor(raw); err == nil {
			st.OutlineColor = c
			note()
		}
	case "outline-offset":
		if l, err := parseLengthAt(raw, *ctx); err == nil && !l.Auto() && !l.None() {
			st.OutlineOffset = l.Resolve(*ctx)
			note()
		}
	case "filter":
		if f, ok := parseFilters(raw, *ctx); ok {
			st.Filters = f
			note()
		}
	case "backdrop-filter":
		if f, ok := parseFilters(raw, *ctx); ok {
			st.BackdropFilters = f
			note()
		}
	case "transform":
		if f, ok := parseTransform(raw, *ctx); ok {
			st.Transform = f
			note()
		}
	case "transform-origin":
		if o, ok := parseTransformOrigin(raw, *ctx); ok {
			st.TransformOrigin = o
			note()
		}
	case "transition":
		if list, ok := parseTransitions(raw); ok {
			st.TransitionProps, st.TransitionDurs, st.TransitionTims, st.TransitionDels = transitionParts(list)
			st.TransitionNone = len(list) == 0
			note()
		}
	case "transition-property":
		if list, none := parsePropertyList(raw); list != nil || none {
			st.TransitionProps = list
			st.TransitionNone = none
			note()
		}
	case "transition-duration":
		if list, ok := parseTimeList(raw); ok {
			st.TransitionDurs = list
			note()
		}
	case "transition-timing-function":
		if list, ok := parseTimingList(raw); ok {
			st.TransitionTims = list
			note()
		}
	case "transition-delay":
		if list, ok := parseTimeList(raw); ok {
			st.TransitionDels = list
			note()
		}
	case "animation":
		if list, ok := parseAnimations(raw); ok {
			st.AnimationNames, st.AnimationDurs, st.AnimationTims, st.AnimationDels, st.AnimationIters, st.AnimationDirs, st.AnimationFills = animationParts(list)
			note()
		}
	case "animation-name":
		if list, ok := parseNameList(raw); ok {
			st.AnimationNames = list
			note()
		}
	case "animation-duration":
		if list, ok := parseTimeList(raw); ok {
			st.AnimationDurs = list
			note()
		}
	case "animation-timing-function":
		if list, ok := parseTimingList(raw); ok {
			st.AnimationTims = list
			note()
		}
	case "animation-delay":
		if list, ok := parseTimeList(raw); ok {
			st.AnimationDels = list
			note()
		}
	case "animation-iteration-count":
		if list, ok := parseIterationList(raw); ok {
			st.AnimationIters = list
			note()
		}
	case "animation-direction":
		if list, ok := parseDirectionList(raw); ok {
			st.AnimationDirs = list
			note()
		}
	case "animation-fill-mode":
		if list, ok := parseFillList(raw); ok {
			st.AnimationFills = list
			note()
		}
	default:
		// Unknown property: ignored, like CSS, and reported so the style
		// author hears what the engine did not use.
		return []string{fmtErrf("ignoring unknown property %q", prop).Error()}
	}
	return nil
}

func setLenSide(sides []Length, canonical string, set map[string]bool, i int, raw string, ctx Units) {
	l, err := parseLengthAt(raw, ctx)
	if err != nil {
		return
	}
	sides[i] = l
	set[canonical] = true
}

// applyBorder sets one or all sides of the border from a shorthand like
// "1px solid red". sides[i] >= 0 names the side to set, or all four.
func applyBorder(st *Style, set map[string]bool, sides [4]int, raw string, note func()) {
	words := splitWords(raw)
	if len(words) == 0 {
		return
	}
	var (
		w                [4]int
		s                [4]uint8
		c                [4]canvas.Color
		setW, setS, setC bool
	)
	for _, wrd := range words {
		if l, err := parseLength(wrd); err == nil && l.u == unitPx && !l.IsPct() {
			for i := 0; i < 4; i++ {
				if sides[i] >= 0 {
					w[sides[i]] = max(l.Px(0), 0)
				} else {
					w[i] = max(l.Px(0), 0)
				}
			}
			setW = true
			continue
		}
		if v, ok := parseBorderStyle(wrd); ok {
			for i := 0; i < 4; i++ {
				if sides[i] >= 0 {
					s[sides[i]] = v
				} else {
					s[i] = v
				}
			}
			setS = true
			continue
		}
		if col, err := ParseColor(wrd); err == nil {
			for i := 0; i < 4; i++ {
				if sides[i] >= 0 {
					c[sides[i]] = col
				} else {
					c[i] = col
				}
			}
			setC = true
		}
	}
	for i := 0; i < 4; i++ {
		if sides[i] >= 0 {
			if setW {
				st.BorderWidth[sides[i]] = w[sides[i]]
			}
			if setS {
				st.BoxStyle[sides[i]] = s[sides[i]]
			}
			if setC {
				st.BoxColor[sides[i]] = c[sides[i]]
			}
		} else {
			if setW {
				st.BorderWidth[i] = w[i]
			}
			if setS {
				st.BoxStyle[i] = s[i]
			}
			if setC {
				st.BoxColor[i] = c[i]
			}
		}
	}
	if setW {
		set["border-width"] = true
	}
	if setS {
		set["border-style"] = true
	}
	if setC {
		set["border-color"] = true
	}
	note()
}

// backgroundKeys are the canonical properties a background shorthand touches,
// marked in the Set map the way a longhand marks only itself.
var backgroundKeys = []string{
	"background", "background-color", "background-image", "background-position",
	"background-size", "background-repeat", "background-clip", "background-origin",
	"background-attachment",
}
