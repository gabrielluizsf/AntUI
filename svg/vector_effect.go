package svg

// VectorEffect is `vector-effect`: which of the transformations the host puts
// on an element the element's own user coordinate system keeps. Zero is the
// drawing's ordinary behaviour, where everything done to the element happens to
// its shape as well.
//
// Every value the spec names is read and given a meaning here rather than
// refused, and the contract they are applied under is this: the host is the
// canvas the drawing comes out on — the furthest SVG viewport, which for a
// drawing painted on its own is its own — so the viewBox fit is part of the
// transform from the element to the host like any other, and a size that does
// not scale is held in the pixels the drawing comes out at rather than in the
// units the file wrote it in. See [vectorEffectMatrix].
//
// The flags do not travel: the spec has the property as not inherited, so
// [Style.with] clears the whole of it before reading what the element itself
// wrote. What does reach the shapes inside an element is the coordinate system
// the effect makes of it, which is why the effect is spent where the picture is
// painted rather than on the element's style.
type VectorEffect struct {
	// NonScalingStroke is the stroke drawn at the width the drawing wrote, in
	// the pixels the drawing comes out at, rather than taken through the
	// viewBox and the transforms that carry it there.
	NonScalingStroke bool
	// NonScalingSize is the element's user coordinate system drawn at its own
	// scale whatever the host scaled it by. The spec has it carry
	// non-scaling-stroke with it, since a scale that does not scale has no
	// width to stretch a stroke with either.
	NonScalingSize bool
	// NonRotation is the element's user coordinate system turned upright: its
	// rotation and skew dropped, its scale kept.
	NonRotation bool
	// FixedPosition is the element's user coordinate system placed where the
	// element's own transform writes rather than where the host moved it to.
	// The host's rotation and scale still happen to it.
	FixedPosition bool
	// OriginX and OriginY are the shift the element's own transform writes, in
	// the element's own coordinates. fixed-position spends it as the origin of
	// the fixed user coordinate system. They are read for every element and
	// mean anything only to one asking for that.
	OriginX, OriginY float64
}

// clear takes the flags off and keeps the origin the element's own transform
// wrote: `none` is the value with nothing in it, and every element starts from
// it before the one it wrote is read.
func (ve *VectorEffect) clear() {
	ve.NonScalingStroke, ve.NonScalingSize, ve.NonRotation, ve.FixedPosition = false, false, false, false
}

// constraining says whether the element asked for a coordinate system that
// resists the host at all, as opposed to only for a stroke that does not scale:
// the first needs the transform the picture is painted through rewritten, and
// the second is settled where the stroke is widened.
func (ve VectorEffect) constraining() bool {
	return ve.NonScalingSize || ve.NonRotation || ve.FixedPosition
}
