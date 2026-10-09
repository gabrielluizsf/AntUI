package css

// Typography constants and parsers. The engine models the text properties
// CSS passes down: weight, style, line-height, spacing, transform and
// decoration, white-space, overflow-wrap, text-overflow and vertical-align.

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
