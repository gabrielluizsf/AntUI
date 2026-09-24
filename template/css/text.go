package css

// Typography constants and parsers. The engine models the text properties
// CSS passes down: weight, style, line-height, spacing, transform and
// decoration, white-space, overflow-wrap, text-overflow and vertical-align.

import (
	"strconv"
	"strings"
)

// FontStyle* are the font-style keywords.
const (
	FontStyleNormal uint8 = iota
	FontStyleItalic
	FontStyleOblique
)

// TextTransform* are the text-transform keywords.
const (
	TextTransformNone uint8 = iota
	TextTransformUppercase
	TextTransformLowercase
	TextTransformCapitalize
)

// TextDecorationUnderline and friends are the bits of text-decoration-line.
const (
	TextDecorationUnderline uint8 = 1 << iota
	TextDecorationOverline
	TextDecorationLineThrough
)

// WhiteSpace* are the white-space keywords.
const (
	WhiteSpaceNormal uint8 = iota
	WhiteSpaceNowrap
	WhiteSpacePre
	WhiteSpacePreWrap
	WhiteSpacePreLine
)

// OverflowWrap* are the overflow-wrap (and word-wrap) keywords.
const (
	OverflowWrapNormal uint8 = iota
	OverflowWrapBreakWord
	OverflowWrapAnywhere
)

// TextOverflow* are the text-overflow keywords.
const (
	TextOverflowClip uint8 = iota
	TextOverflowEllipsis
)

// VerticalAlign* are the vertical-align keywords. A length is stored apart,
// in Style.BaselineShift.
const (
	VerticalAlignBaseline uint8 = iota
	VerticalAlignSub
	VerticalAlignSuper
	VerticalAlignMiddle
	VerticalAlignTop
	VerticalAlignBottom
	VerticalAlignTextTop
	VerticalAlignTextBottom
)

// The absolute weights normal and bold name; the template draws anything at
// 600 or above as bold.
const (
	FontWeightNormal uint16 = 400
	FontWeightBold   uint16 = 700
)

// parseFontWeight reads a font-weight: a number rounded to the nearest
// hundred, or normal/bold/bolder/lighter. The relative keywords cannot see the
// inherited weight in a single declaration, so bolder resolves to bold and
// lighter to 300, which is the common approximation.
func parseFontWeight(raw string) (uint16, bool) {
	s := strings.TrimSpace(raw)
	switch s {
	case "normal":
		return FontWeightNormal, true
	case "bold", "bolder":
		return FontWeightBold, true
	case "lighter":
		return 300, true
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 || v > 1000 {
		return 0, false
	}
	return uint16((v + 50) / 100 * 100), true
}

// parseLineHeight reads a unitless multiplier, a percentage, or a length,
// always computing a multiplier of the font size. Zero means normal.
func parseLineHeight(raw string, ctx Units) (float64, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, false
	}
	if s == "normal" {
		return 0, true
	}
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil || v < 0 {
			return 0, false
		}
		return v / 100, true
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		if v < 0 {
			return 0, false
		}
		return v, true
	}
	l, err := parseLengthAt(s, ctx)
	if err != nil {
		return 0, false
	}
	return float64(l.Resolve(ctx)) / float64(ctx.font()), true
}

// parseTextDecoration reads the text-decoration shorthand or its -line
// longhand: none clears the flags, otherwise underline, overline and
// line-through accumulate. Style and colour words the engine ignores are
// skipped, not rejected.
func parseTextDecoration(raw string) (uint8, bool) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return 0, false
	}
	var bits uint8
	for _, p := range parts {
		switch p {
		case "none":
			return 0, true
		case "underline":
			bits |= TextDecorationUnderline
		case "overline":
			bits |= TextDecorationOverline
		case "line-through":
			bits |= TextDecorationLineThrough
		default:
			// Thickness, style and colour words are not modelled; skip them.
		}
	}
	return bits, true
}

// parseVerticalAlign reads a vertical-align keyword or length. A length lands
// in the returned BaselineShift; keywords leave it zero.
func parseVerticalAlign(raw string, ctx Units) (uint8, Length, bool) {
	switch raw {
	case "baseline":
		return VerticalAlignBaseline, Length{}, true
	case "sub":
		return VerticalAlignSub, Length{}, true
	case "super":
		return VerticalAlignSuper, Length{}, true
	case "middle":
		return VerticalAlignMiddle, Length{}, true
	case "top":
		return VerticalAlignTop, Length{}, true
	case "bottom":
		return VerticalAlignBottom, Length{}, true
	case "text-top":
		return VerticalAlignTextTop, Length{}, true
	case "text-bottom":
		return VerticalAlignTextBottom, Length{}, true
	}
	if l, err := parseLengthAt(raw, ctx); err == nil {
		return VerticalAlignBaseline, l, true
	}
	return 0, Length{}, false
}

// parseFontStyle reads the font-style keyword.
func parseFontStyle(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return FontStyleNormal, true
	case "italic":
		return FontStyleItalic, true
	case "oblique":
		return FontStyleOblique, true
	}
	return 0, false
}

// parseTextTransform reads the text-transform keyword.
func parseTextTransform(raw string) (uint8, bool) {
	switch raw {
	case "none":
		return TextTransformNone, true
	case "uppercase":
		return TextTransformUppercase, true
	case "lowercase":
		return TextTransformLowercase, true
	case "capitalize":
		return TextTransformCapitalize, true
	}
	return 0, false
}

// parseWhiteSpace reads the white-space keyword.
func parseWhiteSpace(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return WhiteSpaceNormal, true
	case "nowrap":
		return WhiteSpaceNowrap, true
	case "pre":
		return WhiteSpacePre, true
	case "pre-wrap":
		return WhiteSpacePreWrap, true
	case "pre-line":
		return WhiteSpacePreLine, true
	}
	return 0, false
}

// parseOverflowWrap reads the overflow-wrap (or word-wrap) keyword.
func parseOverflowWrap(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return OverflowWrapNormal, true
	case "break-word":
		return OverflowWrapBreakWord, true
	case "anywhere":
		return OverflowWrapAnywhere, true
	}
	return 0, false
}

// parseTextOverflow reads the text-overflow keyword. The two-value form
// (clip ellipsis) is accepted with the last keyword winning.
func parseTextOverflow(raw string) (uint8, bool) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return 0, false
	}
	v := TextOverflowClip
	for _, p := range parts {
		switch p {
		case "clip":
			v = TextOverflowClip
		case "ellipsis":
			v = TextOverflowEllipsis
		default:
			return 0, false
		}
	}
	return v, true
}
