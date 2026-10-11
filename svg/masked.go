package svg

// masked is the mask one node is painted under. It hangs off the node rather
// than off the style, because a mask is not inherited — a group is masked once
// as a whole, when its picture is made, rather than masked again at every
// child as it is drawn — which is the same as a clip. See the `mask` field of
// [Node].
type masked struct {
	// def is the mask that was named, which is [maskNowhere] where the drawing
	// could not follow the reference: that keeps nothing, rather than keeping
	// everything, so that a drawing which asked for a mask it does not have
	// shows the element nowhere rather than showing it whole.
	def *maskDef
	// x, y, w, h is where the mask reaches, in the drawing's own coordinates,
	// and cuts says the region narrows it at all: a region measured against a
	// box that cannot be measured leaves the mask reaching everywhere. See
	// [maskUnder].
	x, y, w, h float64
	cuts       bool
	// box is the box of the element this mask was put on, kept for the content
	// written in fractions of it: [paintMasked] lays that content out over the
	// box, and the box is measured only once, here. It is the same box the
	// region above is a fraction of, so the content and the region reach the
	// same places. See [maskUnder].
	box [4]float64
}
