package svg

import "github.com/gabrielluizsf/antui/canvas"

// clipRegion is the clip one element carries: the shapes of the clipPath it
// named, with its own transform already put on them so that the cut follows
// the element it cuts. It hangs off the node rather than off the style because
// a clip is not inherited (see the `clip` field of [Node]), and it can be
// empty, which is a clip that keeps nothing rather than no clip at all.
type clipRegion struct {
	// boxUnits is whether the shapes were written as fractions of the box of
	// the element being clipped rather than in the drawing's own coordinates.
	boxUnits bool
	shapes   []canvas.MaskShape
	// pending is the id of a clipPath whose shapes are not worked out yet: a
	// `clip-path` written inside a `<clipPath>`, named while every clipPath in
	// the drawing was being built and before any of them was measured. It
	// carries the name alone and no shapes; [Image.resolveClipsIn] follows it
	// when the picture it was written in is walked for its own shapes. See
	// [Image.resolveClip].
	pending string
}
