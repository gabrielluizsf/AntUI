package svg

import "github.com/gabrielluizsf/antui/canvas"

// clipPath is one `<clipPath>` of a drawing: the picture it held, the shapes
// that picture came out as and how it said they were to be measured.
type clipPath struct {
	// shapes are what the clip is made of, each with the rule its own
	// `clip-rule` asked for. The clip is the union of them, so every shape is
	// measured on its own and what any one of them covers is inside — two
	// shapes overlapping leave their overlap standing rather than cutting it
	// away again.
	shapes []canvas.MaskShape
	// boxUnits is whether the shapes were written as fractions of the box of
	// the element being clipped rather than in the drawing's own coordinates,
	// which is the second of the two things `clipPathUnits` may say. See
	// [Image.resolveClip].
	boxUnits bool
	// content is the tree the shapes are read from, kept because a `clip-path`
	// written inside a `<clipPath>` cuts the shape it is written on, and the
	// shapes of the outer one can only be worked out once every clipPath in
	// the drawing is known — one of them may be named further down the file
	// than the place it is pointed at. See [Image.readClipPaths].
	content *Node
	// built says the shapes above are finished: the picture in content was
	// walked and every `clip-path` inside it was followed into them. A
	// reference met while this is still false is one from inside a picture
	// still being read, and waits to be followed — see [Image.resolveClip]
	// and [Image.resolveClipsIn].
	built bool
}
