package css

// textStyle carries the typography fields of a Style, measured in reference
// pixels at scale 1. TextAlign 0 left, 1 center, 2 right, 3 justify. The
// typography group is inherited.
type textStyle struct {
	FontSize       int // 0 means the theme default
	TextAlign      uint8
	FontWeight     uint16  // 0 means normal (400)
	FontStyle      uint8   // one of the FontStyle* constants
	FontFamily     string  // the font-family list, lowercased; empty means none was asked for
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
}
