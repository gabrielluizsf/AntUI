package css

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// Display is the CSS outer display type of a box. Inline and inline-block
// boxes flow along a line, the rest start a new line.
const (
	DisplayBlock uint8 = iota
	DisplayNone
	DisplayInline
	DisplayInlineBlock
	DisplayFlex
	DisplayGrid
	DisplayTable
)

// Position is the CSS positioning scheme. Absolute and fixed boxes are taken
// out of flow; relative and sticky boxes keep their flow slot and only shift
// where they are painted.
const (
	PositionStatic uint8 = iota
	PositionAbsolute
	PositionFixed
	PositionRelative
	PositionSticky
)

// Overflow is the per-axis overflow behaviour of a box.
const (
	OverflowVisible uint8 = iota
	OverflowHidden
	OverflowScroll
	OverflowAuto
	OverflowClip
)

// Visibility controls whether a box is painted; a hidden box keeps its slot.
const (
	VisibilityVisible uint8 = iota
	VisibilityHidden
	VisibilityCollapse
)

// PointerEvents says whether a box takes part in hit testing.
const (
	PointerEventsAuto uint8 = iota
	PointerEventsNone
)

// Cursor names the pointer shape a box asks for. The engine only models it;
// painting a real OS cursor is the window's job.
const (
	CursorDefault uint8 = iota
	CursorPointer
	CursorText
	CursorWait
	CursorCrosshair
	CursorMove
	CursorNotAllowed
	CursorGrab
	CursorGrabbing
	CursorCell
)

// Border styles, matching the values stored in Style.BoxStyle.
const (
	BorderNone uint8 = iota
	BorderSolid
	BorderDashed
	BorderDotted
	BorderDouble
)

// Style is the computed drawing style of one widget after the cascade: which
// pixels to fill, how thick its border is, how big its text. Zero values fall
// back to the template's theme — a style does not know the theme, it knows
// only what the stylesheet said, and the Set map records whether it spoke.
type Style struct {
	// Box model.
	Width, Height        Length
	MinWidth, MaxWidth   Length
	MinHeight, MaxHeight Length
	Margin               [4]Length // top, right, bottom, left
	Padding              [4]Length
	BoxSizing            bool // border-box when true, content-box otherwise

	// Border. BoxStyle[i] is one of the Border* constants.
	BorderWidth [4]int
	BoxStyle    [4]uint8
	BoxColor    [4]canvas.Color
	Radius      [4]int
	RadiusY     [4]int // vertical radii; zero means the same as Radius

	// Surface and ink.
	Background canvas.Color
	// Background layers, painted above the colour in order. The parallel
	// lists cycle over the layers, so one background-position value
	// positions an entire stack of images.
	BackgroundImages []BackImage
	BackgroundPos    []BackPos
	BackgroundSize   []BackSize
	BackgroundRepeat []BackRepeat
	BackgroundClip   []uint8
	BackgroundOrigin []uint8
	BackgroundAttach []uint8
	Opacity          float64 // 0..1; 1 when unset
	Color            canvas.Color

	// Text measured in reference pixels at scale 1. TextAlign 0 left, 1
	// center, 2 right, 3 justify. The typography group is inherited.
	FontSize       int // 0 means the theme default
	TextAlign      uint8
	FontWeight     uint16  // 0 means normal (400)
	FontStyle      uint8   // one of the FontStyle* constants
	LineHeight     float64 // multiplier of the font size; 0 means normal
	LetterSpacing  Length
	WordSpacing    Length
	TextTransform  uint8 // one of the TextTransform* constants
	TextDecoration uint8 // bit set of the TextDecoration* constants
	WhiteSpace     uint8 // one of the WhiteSpace* constants
	OverflowWrap   uint8 // one of the OverflowWrap* constants
	TextOverflow   uint8 // one of the TextOverflow* constants
	VerticalAlign  uint8 // one of the VerticalAlign* constants
	BaselineShift  Length

	// Layout.
	Display       uint8 // one of the Display* constants
	Position      uint8 // one of the Position* constants
	Top           Length
	Right         Length
	Bottom        Length
	Left          Length
	ZIndex        int
	Overflow      [2]uint8 // x, y; one of the Overflow* constants
	Visibility    uint8    // one of the Visibility* constants
	PointerEvents uint8    // one of the PointerEvents* constants
	Cursor        uint8    // one of the Cursor* constants

	// Flexbox container geometry: how a DisplayFlex box lays its children
	// out along the main axis and across it. RowGap and ColumnGap space the
	// items and the wrapped lines.
	FlexDirection  uint8 // one of the FlexDirection* constants
	FlexWrap       uint8 // one of the FlexWrap* constants
	JustifyContent uint8 // one of the Justify* constants
	AlignItems     uint8 // one of the Align* constants
	AlignContent   uint8 // one of the Content* constants
	RowGap         Length
	ColumnGap      Length

	// Flex item geometry: how this box answers its flex container. Order
	// reorders items; grow and shrink share the free space; basis is the
	// item's main size before distribution; align-self overrides the
	// container's align-items.
	Order      int
	FlexGrow   float64
	FlexShrink float64
	FlexBasis  Length
	AlignSelf  uint8 // AlignAuto or one of the Align* constants

	// Visual effects. OutlineStyle shares the Border* constants.
	BoxShadow       []Shadow
	TextShadow      []Shadow
	OutlineWidth    int
	OutlineStyle    uint8
	OutlineColor    canvas.Color
	OutlineOffset   int
	Filters         []Filter
	BackdropFilters []Filter

	// Transform is the transform list, applied around the border box's
	// TransformOrigin. Nil means none.
	Transform []TransformFunc
	// TransformOrigin is the transform pivot: the initial value sits in the
	// box's centre.
	TransformOrigin [2]Length

	// Transitions is the finished transition list: which properties animate
	// toward new computed values, and how, the moment the cascade changes
	// them. Animations is the finished animation list, the keyframes blocks
	// the element plays and how. Both are assembled from the parallel
	// longhand lists below by finishTransitions and finishAnimations at the
	// end of the cascade; the longhand lists are intermediate and never
	// meant for the template.
	Transitions []Transition
	Animations  []Animation

	TransitionProps []string
	TransitionDurs  []Time
	TransitionTims  []Timing
	TransitionDels  []Time
	TransitionNone  bool

	AnimationNames []string
	AnimationDurs  []Time
	AnimationTims  []Timing
	AnimationDels  []Time
	AnimationIters []float64 // math.Inf(+1) means infinite
	AnimationDirs  []uint8
	AnimationFills []uint8

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

// BorderOn reports whether any side paints a border.
func (s *Style) BorderOn() bool {
	return s.BoxStyle[0] != BorderNone || s.BoxStyle[1] != BorderNone ||
		s.BoxStyle[2] != BorderNone || s.BoxStyle[3] != BorderNone
}

// Inline reports whether the box flows along a line with its neighbours.
func (s Style) Inline() bool {
	return s.Display == DisplayInline || s.Display == DisplayInlineBlock
}

// OutOfFlow reports whether the box is positioned outside the normal flow.
func (s Style) OutOfFlow() bool {
	return s.Position == PositionAbsolute || s.Position == PositionFixed
}

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

// resolveVars replaces every var(--name[, fallback]) in a raw value with the
// cascaded custom-property value, or the fallback, walking nested definitions
// up to a depth guard so a cycle --a: var(--a) cannot spin forever. A var()
// with no known value and no fallback is left in place, which makes the
// property unrecognised and dropped by applyDecl.
func resolveVars(raw string, customs map[string]string) string {
	return resolveVarsDepth(raw, customs, 0)
}

func resolveVarsDepth(raw string, customs map[string]string, depth int) string {
	if depth > 16 {
		return raw
	}
	out := raw
	for i := 0; i < len(out); i++ {
		if out[i] != 'v' || !strings.HasPrefix(out[i:], "var(") {
			continue
		}
		end := parenEnd(out, i+3)
		if end < 0 {
			return out
		}
		name, fallback, ok := splitVarArg(out[i+4 : end])
		if !ok {
			i = end
			continue
		}
		if v, found := customs[name]; found {
			out = out[:i] + resolveVarsDepth(v, customs, depth+1) + out[end+1:]
			continue
		}
		if fallback != "" {
			out = out[:i] + resolveVarsDepth(fallback, customs, depth+1) + out[end+1:]
			continue
		}
		i = end
	}
	return out
}

// parenEnd returns the index of the ')' matching the '(' at open, or -1.
func parenEnd(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		case '"', '\'':
			if j := quoteEnd(s, i); j >= 0 {
				i = j
			}
		}
	}
	return -1
}

// quoteEnd returns the index of the quote closing the one at i, or -1.
func quoteEnd(s string, i int) int {
	q := s[i]
	for j := i + 1; j < len(s); j++ {
		if s[j] == '\\' {
			j++
			continue
		}
		if s[j] == q {
			return j
		}
	}
	return -1
}

// splitVarArg reads " --name , fallback " into the custom-property name, the
// fallback (may be empty) and ok. The fallback is everything after the first
// top-level comma.
func splitVarArg(inner string) (name, fallback string, ok bool) {
	inside := false
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		switch c {
		case '"', '\'':
			if j := quoteEnd(inner, i); j >= 0 {
				i = j
			}
		case '(', '[':
			inside = true
		case ')', ']':
			inside = false
		case ',':
			if !inside {
				name = strings.TrimSpace(inner[:i])
				fallback = strings.TrimSpace(inner[i+1:])
				return name, fallback, name != ""
			}
		}
	}
	return strings.TrimSpace(inner), "", strings.TrimSpace(inner) != ""
}

func setLenSide(sides []Length, canonical string, set map[string]bool, i int, raw string, ctx Units) {
	l, err := parseLengthAt(raw, ctx)
	if err != nil {
		return
	}
	sides[i] = l
	set[canonical] = true
}

func parseOpacity(raw string) (float64, bool) {
	var v float64
	if _, err := scanFloat(raw, &v); err != nil {
		return 0, false
	}
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return v, true
}

func parseBorderWidths(raw string) ([4]int, bool) {
	var out [4]int
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]int, len(parts))
	for i, p := range parts {
		l, err := parseLength(p)
		if err != nil || l.u != unitPx || l.value < 0 {
			return out, false
		}
		vals[i] = int(l.value + 0.5)
	}
	fill := expandFour(vals)
	copy(out[:], fill[:])
	return out, true
}

func parseBorderStyles(raw string) ([4]uint8, bool) {
	var out [4]uint8
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]uint8, len(parts))
	for i, p := range parts {
		v, ok := parseBorderStyle(p)
		if !ok {
			return out, false
		}
		vals[i] = v
	}
	fill := expandFourInt8(vals)
	copy(out[:], fill[:])
	return out, true
}

func parseDisplay(raw string) (uint8, bool) {
	switch raw {
	case "block", "flow-root", "list-item":
		return DisplayBlock, true
	case "none":
		return DisplayNone, true
	case "inline":
		return DisplayInline, true
	case "inline-block", "inline-flex", "inline-grid", "inline-table":
		return DisplayInlineBlock, true
	case "flex":
		return DisplayFlex, true
	case "grid":
		return DisplayGrid, true
	case "table", "table-row", "table-cell", "table-caption", "table-column",
		"table-column-group", "table-header-group", "table-footer-group",
		"table-row-group":
		return DisplayTable, true
	case "contents":
		return DisplayNone, true
	}
	return 0, false
}

func parsePosition(raw string) (uint8, bool) {
	switch raw {
	case "static":
		return PositionStatic, true
	case "absolute":
		return PositionAbsolute, true
	case "fixed":
		return PositionFixed, true
	case "relative":
		return PositionRelative, true
	case "sticky":
		return PositionSticky, true
	}
	return 0, false
}

func parseOverflow(raw string) (uint8, bool) {
	switch raw {
	case "visible":
		return OverflowVisible, true
	case "hidden":
		return OverflowHidden, true
	case "scroll":
		return OverflowScroll, true
	case "auto":
		return OverflowAuto, true
	case "clip":
		return OverflowClip, true
	}
	return 0, false
}

func parseVisibility(raw string) (uint8, bool) {
	switch raw {
	case "visible":
		return VisibilityVisible, true
	case "hidden":
		return VisibilityHidden, true
	case "collapse":
		return VisibilityCollapse, true
	}
	return 0, false
}

func parsePointerEvents(raw string) (uint8, bool) {
	switch raw {
	case "auto":
		return PointerEventsAuto, true
	case "none":
		return PointerEventsNone, true
	}
	return 0, false
}

func parseCursor(raw string) (uint8, bool) {
	switch raw {
	case "default", "auto":
		return CursorDefault, true
	case "pointer":
		return CursorPointer, true
	case "text":
		return CursorText, true
	case "wait", "progress":
		return CursorWait, true
	case "crosshair":
		return CursorCrosshair, true
	case "move", "all-scroll":
		return CursorMove, true
	case "not-allowed":
		return CursorNotAllowed, true
	case "grab":
		return CursorGrab, true
	case "grabbing":
		return CursorGrabbing, true
	case "cell":
		return CursorCell, true
	}
	return 0, false
}

// parseZIndex accepts an integer (possibly negative) z-index. auto is the
// initial value and resets the layer to 0.
func parseZIndex(raw string) (int, bool) {
	if raw == "auto" {
		return 0, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}

func parseBorderColors(raw string) ([4]canvas.Color, bool) {
	var out [4]canvas.Color
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]canvas.Color, len(parts))
	for i, p := range parts {
		c, err := ParseColor(p)
		if err != nil {
			return out, false
		}
		vals[i] = c
	}
	fill := expandFourColor(vals)
	copy(out[:], fill[:])
	return out, true
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

func parseBorderStyle(raw string) (uint8, bool) {
	switch raw {
	case "solid":
		return BorderSolid, true
	case "dashed":
		return BorderDashed, true
	case "dotted":
		return BorderDotted, true
	case "double":
		return BorderDouble, true
	case "groove", "ridge", "inset", "outset":
		// Drawn as a solid ring: a flat canvas has no bevel to shade.
		return BorderSolid, true
	case "none", "hidden":
		return BorderNone, true
	}
	return 0, false
}

// parseRadiiList expands one to four non-negative pixel radii into the four
// corners, in CSS order: top-left, top-right, bottom-right, bottom-left.
func parseRadiiList(raw string) ([4]int, bool) {
	var out [4]int
	parts := splitWords(raw)
	if len(parts) == 0 || len(parts) > 4 {
		return out, false
	}
	vals := make([]int, len(parts))
	for i, p := range parts {
		l, err := parseLength(p)
		if err != nil || l.u != unitPx || l.value < 0 {
			return out, false
		}
		vals[i] = int(l.value + 0.5)
	}
	fill := expandFour(vals)
	copy(out[:], fill[:])
	return out, true
}

// parseRadiiXY parses a border-radius value, splitting the elliptical form
// "8px / 16px" into horizontal and vertical radii. With no slash the two axes
// are equal.
func parseRadiiXY(raw string) ([4]int, [4]int, bool) {
	head, tail, hasSlash := splitSlash(raw)
	rx, ok := parseRadiiList(head)
	if !ok {
		return rx, rx, false
	}
	ry := rx
	if hasSlash && strings.TrimSpace(tail) != "" {
		if v, ok := parseRadiiList(tail); ok {
			ry = v
		} else {
			return rx, rx, false
		}
	}
	return rx, ry, true
}

func expandFour(v []int) [4]int {
	switch len(v) {
	case 1:
		return [4]int{v[0], v[0], v[0], v[0]}
	case 2:
		return [4]int{v[0], v[1], v[0], v[1]}
	case 3:
		return [4]int{v[0], v[1], v[2], v[1]}
	}
	return [4]int{v[0], v[1], v[2], v[3]}
}

func expandFourInt8(v []uint8) [4]uint8 {
	switch len(v) {
	case 1:
		return [4]uint8{v[0], v[0], v[0], v[0]}
	case 2:
		return [4]uint8{v[0], v[1], v[0], v[1]}
	case 3:
		return [4]uint8{v[0], v[1], v[2], v[1]}
	}
	return [4]uint8{v[0], v[1], v[2], v[3]}
}

func expandFourColor(v []canvas.Color) [4]canvas.Color {
	switch len(v) {
	case 1:
		return [4]canvas.Color{v[0], v[0], v[0], v[0]}
	case 2:
		return [4]canvas.Color{v[0], v[1], v[0], v[1]}
	case 3:
		return [4]canvas.Color{v[0], v[1], v[2], v[1]}
	}
	return [4]canvas.Color{v[0], v[1], v[2], v[3]}
}

// inheritedProps names the properties CSS passes down the tree by default:
// the textual ones. A template's widgets have no nested parent, so the
// inheritance source is the body's computed style — the box properties are
// deliberately absent, they do not inherit in CSS either. Every property
// implemented later that inherits by default must be added here.
var inheritedProps = map[string]bool{
	"color":          true,
	"font-size":      true,
	"font-weight":    true,
	"font-style":     true,
	"line-height":    true,
	"letter-spacing": true,
	"word-spacing":   true,
	"text-align":     true,
	"text-transform": true,
	"white-space":    true,
	"overflow-wrap":  true,
	"visibility":     true,
	"cursor":         true,
	"text-shadow":    true,
}

// backgroundKeys are the canonical properties a background shorthand touches,
// marked in the Set map the way a longhand marks only itself.
var backgroundKeys = []string{
	"background", "background-color", "background-image", "background-position",
	"background-size", "background-repeat", "background-clip", "background-origin",
	"background-attachment",
}

// applyKeyword resolves one of the four cascade keywords on a property.
// initial and revert reset the property to its CSS initial value; inherit (and
// unset on an inherited property) record the intent for StyleUnits to fold in
// the body's computed value after the cascade. note() keeps the cascade
// honest: a real value that wins later clears the recorded intent, and so does
// a later keyword.
func applyKeyword(st *Style, note func(), kw, prop string) []string {
	switch kw {
	case "initial", "revert":
		applyInitial(st, prop)
		note()
		return nil
	case "inherit":
		note()
		st.inherit[prop] = true
		return nil
	case "unset":
		if inheritedProps[prop] {
			note()
			st.inherit[prop] = true
			return nil
		}
		applyInitial(st, prop)
		note()
		return nil
	}
	return nil
}

// ThemeInk is what the engine stores in Style.Color when the cascade resolves
// colour to its CSS initial value — a `color: initial` (or revert/unset-less
// body) means the theme's ink, not a literal pixel. Templates must translate
// this sentinel into the window's theme text colour; the sentinel itself is an
// almost-invisible black that must never be produced by parseColor.
const ThemeInk = canvas.Color(0xFF000001)

// applyInitial resets a property to the CSS initial value. Where the initial
// and the engine's "unset" carry the same shape — a zero colour meaning the
// theme's ink, width auto — they stay the same field value apart from colour,
// which keeps its own sentinel so templates can tell "theme ink" apart from a
// declared transparent.
func applyInitial(st *Style, prop string) {
	switch prop {
	case "color":
		st.Color = ThemeInk
	case "opacity":
		st.Opacity = 1
	case "font-size":
		st.FontSize = 0
	case "text-align":
		st.TextAlign = 0
	case "font-weight":
		st.FontWeight = 0
	case "font-style":
		st.FontStyle = FontStyleNormal
	case "line-height":
		st.LineHeight = 0
	case "letter-spacing":
		st.LetterSpacing = Zero()
	case "word-spacing":
		st.WordSpacing = Zero()
	case "text-transform":
		st.TextTransform = TextTransformNone
	case "text-decoration", "text-decoration-line":
		st.TextDecoration = 0
	case "white-space":
		st.WhiteSpace = WhiteSpaceNormal
	case "overflow-wrap", "word-wrap":
		st.OverflowWrap = OverflowWrapNormal
	case "text-overflow":
		st.TextOverflow = TextOverflowClip
	case "vertical-align":
		st.VerticalAlign = VerticalAlignBaseline
		st.BaselineShift = Length{}
	case "width":
		st.Width = Auto()
	case "min-width":
		st.MinWidth = Auto()
	case "height":
		st.Height = Auto()
	case "min-height":
		st.MinHeight = Auto()
	case "max-width":
		st.MaxWidth = Length{u: unitNone}
	case "max-height":
		st.MaxHeight = Length{u: unitNone}
	case "margin":
		st.Margin = [4]Length{Zero(), Zero(), Zero(), Zero()}
	case "margin-top":
		st.Margin[0] = Zero()
	case "margin-right":
		st.Margin[1] = Zero()
	case "margin-bottom":
		st.Margin[2] = Zero()
	case "margin-left":
		st.Margin[3] = Zero()
	case "padding":
		st.Padding = [4]Length{Zero(), Zero(), Zero(), Zero()}
	case "padding-top":
		st.Padding[0] = Zero()
	case "padding-right":
		st.Padding[1] = Zero()
	case "padding-bottom":
		st.Padding[2] = Zero()
	case "padding-left":
		st.Padding[3] = Zero()
	case "box-sizing":
		st.BoxSizing = false
	case "background":
		st.Background = canvas.Transparent
		st.BackgroundImages = nil
		st.BackgroundPos = nil
		st.BackgroundSize = nil
		st.BackgroundRepeat = nil
		st.BackgroundClip = nil
		st.BackgroundOrigin = nil
		st.BackgroundAttach = nil
	case "background-color":
		st.Background = canvas.Transparent
	case "background-image":
		st.BackgroundImages = nil
	case "background-position":
		st.BackgroundPos = nil
	case "background-size":
		st.BackgroundSize = nil
	case "background-repeat":
		st.BackgroundRepeat = nil
	case "background-clip":
		st.BackgroundClip = nil
	case "background-origin":
		st.BackgroundOrigin = nil
	case "background-attachment":
		st.BackgroundAttach = nil
	case "border":
		st.BorderWidth = [4]int{}
		st.BoxStyle = [4]uint8{}
		st.BoxColor = [4]canvas.Color{}
	case "border-width":
		st.BorderWidth = [4]int{}
	case "border-style":
		st.BoxStyle = [4]uint8{}
	case "border-color":
		st.BoxColor = [4]canvas.Color{}
	case "border-radius":
		st.Radius = [4]int{}
		st.RadiusY = [4]int{}
	case "display":
		st.Display = DisplayBlock
	case "position":
		st.Position = PositionStatic
	case "top":
		st.Top = Auto()
	case "right":
		st.Right = Auto()
	case "bottom":
		st.Bottom = Auto()
	case "left":
		st.Left = Auto()
	case "z-index":
		st.ZIndex = 0
	case "overflow":
		st.Overflow = [2]uint8{OverflowVisible, OverflowVisible}
	case "overflow-x":
		st.Overflow[0] = OverflowVisible
	case "overflow-y":
		st.Overflow[1] = OverflowVisible
	case "visibility":
		st.Visibility = VisibilityVisible
	case "pointer-events":
		st.PointerEvents = PointerEventsAuto
	case "cursor":
		st.Cursor = CursorDefault
	case "flex-direction":
		st.FlexDirection = FlexDirectionRow
	case "flex-wrap":
		st.FlexWrap = FlexWrapNowrap
	case "flex-flow":
		st.FlexDirection = FlexDirectionRow
		st.FlexWrap = FlexWrapNowrap
	case "justify-content":
		st.JustifyContent = JustifyFlexStart
	case "align-items":
		st.AlignItems = AlignStretch
	case "align-self":
		st.AlignSelf = AlignAuto
	case "align-content":
		st.AlignContent = ContentStretch
	case "gap":
		st.RowGap = Zero()
		st.ColumnGap = Zero()
	case "row-gap":
		st.RowGap = Zero()
	case "column-gap":
		st.ColumnGap = Zero()
	case "order":
		st.Order = 0
	case "flex":
		st.FlexGrow = 0
		st.FlexShrink = 1
		st.FlexBasis = Auto()
	case "flex-grow":
		st.FlexGrow = 0
	case "flex-shrink":
		st.FlexShrink = 1
	case "flex-basis":
		st.FlexBasis = Auto()
	case "box-shadow":
		st.BoxShadow = nil
	case "text-shadow":
		st.TextShadow = nil
	case "outline":
		st.OutlineWidth = 0
		st.OutlineStyle = BorderNone
		st.OutlineColor = CurrentColor
	case "outline-width":
		st.OutlineWidth = 0
	case "outline-style":
		st.OutlineStyle = BorderNone
	case "outline-color":
		st.OutlineColor = CurrentColor
	case "outline-offset":
		st.OutlineOffset = 0
	case "filter":
		st.Filters = nil
	case "backdrop-filter":
		st.BackdropFilters = nil
	case "transform":
		st.Transform = nil
	case "transform-origin":
		st.TransformOrigin = InitialTransformOrigin
	case "transition":
		st.TransitionProps = nil
		st.TransitionDurs = nil
		st.TransitionTims = nil
		st.TransitionDels = nil
		st.TransitionNone = false
	case "transition-property":
		st.TransitionProps = nil
		st.TransitionNone = false
	case "transition-duration":
		st.TransitionDurs = nil
	case "transition-timing-function":
		st.TransitionTims = nil
	case "transition-delay":
		st.TransitionDels = nil
	case "animation":
		st.AnimationNames = nil
		st.AnimationDurs = nil
		st.AnimationTims = nil
		st.AnimationDels = nil
		st.AnimationIters = nil
		st.AnimationDirs = nil
		st.AnimationFills = nil
	case "animation-name":
		st.AnimationNames = nil
	case "animation-duration":
		st.AnimationDurs = nil
	case "animation-timing-function":
		st.AnimationTims = nil
	case "animation-delay":
		st.AnimationDels = nil
	case "animation-iteration-count":
		st.AnimationIters = nil
	case "animation-direction":
		st.AnimationDirs = nil
	case "animation-fill-mode":
		st.AnimationFills = nil
	}
}
