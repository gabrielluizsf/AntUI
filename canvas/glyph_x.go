package canvas

import "github.com/gabrielluizsf/antui/font"

// glyph returns the 16 rows of the built-in font for a codepoint.
func glyph(r rune) []byte { return font.Glyph(r) }
