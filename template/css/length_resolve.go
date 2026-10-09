package css

import (
	"math"
)

// Resolve turns the length into reference pixels under a context. Keywords
// (auto, none) resolve to zero; a caller that needs to tell them apart asks
// [Length.Auto] or [Length.None] first.
func (l Length) Resolve(ctx Units) int {
	return int(math.Round(l.ref(ctx)))
}

// Px resolves the length against a base: percentages scale the base, fixed
// lengths return their pixels, keywords return zero. It is [Length.Resolve]
// with the base standing in for the whole measurement context, which keeps
// the old one-argument call sites working.
func (l Length) Px(base int) int {
	return l.Resolve(Units{Width: base, Height: base, Font: base, Root: base})
}
