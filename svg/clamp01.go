package svg

import "math"

// clamp01 keeps a number a channel can be: between zero and one, where the
// operations may push it past what a colour holds.
func clamp01(v float64) float64 { return math.Max(0, math.Min(1, v)) }
