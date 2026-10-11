package template

// gridIntrinsicWidth is what the container's column tracks ask for on their
// own: the narrowest, which is the automatic minimum of every track summed, and
// the widest, which is every track at the top of its own range. Gaps count
// once. The first is what a min-content track of a parent wants; the second is
// the width the container would take if the parent asked how wide it would like
// to be.
func (b *gridBatch) gridIntrinsicWidth() (narrowest, widest int) {
	contributions := b.gridTrackContributions(len(b.trackColumns), true)
	for i, track := range b.trackColumns {
		plan := b.gridTrackPlan(track, contributions[i], b.contentW, gridUnbounded, true)
		narrowest += plan.min
		widest += min(plan.max, max(plan.size, plan.min))
	}
	gaps := max(len(b.trackColumns)-1, 0) * b.colGap
	return narrowest + gaps, widest + gaps
}
