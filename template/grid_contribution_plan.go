package template

// gridContributionPlan is a track sized from its content: the size it settles
// on, the narrowest and the widest its items allow, and whether the space left
// over may still be handed to it.
func gridContributionPlan(contribution gridContribution, size int, stretch bool) gridTrackPlan {
	return gridTrackPlan{
		size:    max(size, contribution.min),
		min:     contribution.min,
		max:     max(contribution.max, contribution.min),
		stretch: stretch,
	}
}
