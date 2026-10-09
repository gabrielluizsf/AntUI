package css

// constrained reports whether the query measures the viewport at all. One
// that does not — a media type on its own, or a condition the canvas has no
// sensor for — matches every window, which is what keeps a rule from being
// dropped for want of a measurement.
func (q MediaQuery) constrained() bool {
	return q.HasMin || q.HasMax || q.HasMinHeight || q.HasMaxHeight ||
		q.HasMinDpi || q.HasMaxDpi ||
		q.HasOrientation || q.HasScheme
}
