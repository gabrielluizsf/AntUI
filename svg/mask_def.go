package svg

import "github.com/gabrielluizsf/antui/canvas"

// maskDef is one `<mask>` of a drawing: the picture it holds, how that picture
// says how much it keeps, and where it reaches.
type maskDef struct {
	// content is the mask's own tree, built in the coordinates it was written
	// in — the transform on the mask and on anything inside it counts, since
	// that is part of how they were written — and it is the element that names
	// the mask that moves it afterwards. It is nil only for [maskNowhere], the
	// mask a reference that the drawing cannot follow stands for.
	content *Node
	// measure is how the content says how much of what it covers is kept: by
	// how bright it is, which is what a mask means where nothing has said
	// otherwise, or by how much of it is there, which is what `mask-type:
	// alpha` asks for.
	measure canvas.MaskMode
	// region is where the mask reaches, and boxUnits inside it says which of
	// the two ways `maskUnits` has of writing that it was written in. See
	// [readMaskRegion] and [maskUnder].
	region maskRegion
	// contentBoxUnits is whether the content was written in fractions of the
	// box of the element the mask is put on rather than in the drawing's own
	// coordinates, which is the second of the two things `maskContentUnits`
	// may say. The box it is a fraction of is the one the element paints over,
	// worked out where the element is built. See [maskUnder].
	contentBoxUnits bool
}
