package canvas

// The two ways a mask may measure. Luminance is what an SVG `<mask>` means
// unless it says otherwise: what is bright keeps and what is dark takes away.
// Alpha keeps by how much of the mask is there at all, so an opaque black
// shape keeps everything under it and only what is clear takes away — the
// difference between the two is a black mask, which under one is a hole and
// under the other is no mask at all.
const (
	MaskLuminance MaskMode = iota
	MaskAlpha
)
