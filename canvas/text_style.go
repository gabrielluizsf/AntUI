package canvas

// TextStyle is how a run of text is drawn: its colour, an integer enlargement
// of the face, a fake bold (the stroke drawn twice), a fake italic (each glyph
// row shifted to lean right) and the extra spacing CSS asks for after every
// glyph and after every space.
type TextStyle struct {
	// Color of the drawn text.
	Color Color

	// Face is the face to draw in when the style is not taking the one the
	// canvas already holds — a face a font-family named and this program
	// holds a file for. nil means that face.
	Face          *Face
	Scale         int // 0 or 1 draw the face at its own size
	Bold          bool
	Italic        bool
	LetterSpacing int // extra pixels after every glyph
	WordSpacing   int // extra pixels after every space
}
