package svg

import "strings"

// withVectorEffect reads `vector-effect`: the ways an element's own coordinate
// system is held still against the transformations the host puts on it. The
// spec's values are read alone or together, and a value this package does not
// know is said once and left out rather than the whole declaration being
// refused, so a drawing written against a newer spec keeps the rest of what it
// asked for. `none` is the empty value and clears whatever the list has said so
// far; nothing written is silence, as everywhere else.
//
// `non-scaling-stroke` is the one value that says nothing about the shape, so
// it is kept beside the other three, which rewrite the transform the element
// and everything inside it are painted through — see [vectorEffectMatrix].
func (s *Style) withVectorEffect(raw string, warn func(string, ...any)) {
	s.VectorEffect.clear()
	for _, name := range strings.Fields(strings.ReplaceAll(raw, ",", " ")) {
		switch strings.ToLower(name) {
		case "none":
			s.VectorEffect.clear()
		case "non-scaling-stroke":
			s.VectorEffect.NonScalingStroke = true
		case "non-scaling-size":
			// The spec has this one carry a stroke that does not scale with
			// it: a scale that does not scale has no width to stretch one with.
			s.VectorEffect.NonScalingSize, s.VectorEffect.NonScalingStroke = true, true
		case "non-rotation":
			s.VectorEffect.NonRotation = true
		case "fixed-position":
			s.VectorEffect.FixedPosition = true
		default:
			warn("the vector-effect %q is not a value this package draws, so it is left out", name)
		}
	}
}
