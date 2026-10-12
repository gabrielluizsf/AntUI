package template

import "github.com/gabrielluizsf/antui/template/css"

// gridTrackPlan sizes one track. A track written as a length is exactly that
// length; a fraction grows into what the flexible tracks share; auto takes what
// its items ask for and stretches into the space left over; and the content
// keywords ask for the narrowest or the widest they can be, with fit-content()
// clamped between the two by the limit it was given.
func (b *gridBatch) gridTrackPlan(track css.GridTrack, contribution gridContribution, base, available int, definite bool) gridTrackPlan {
	if track.Kind == css.GridTrackSingle {
		return b.singleGridTrackPlan(track, contribution, base, available, definite)
	}

	minimum := b.gridTrackMin(track.Min, contribution, base, !definite)
	plan := gridTrackPlan{min: minimum}
	switch track.Max.Kind {
	case css.GridTrackLength:
		maximum := contribution.max
		if definite || !track.Max.Length.IsPct() {
			maximum = b.css.style.length(b.st, track.Max.Length, base)
		}
		plan.max = max(maximum, minimum)
		plan.size = min(max(contribution.max, minimum), plan.max)
	case css.GridTrackFlexible:
		plan.max = contribution.max
		if definite {
			plan.max = max(available, minimum)
		}
		plan.max = max(plan.max, minimum)
		plan.size = min(max(contribution.max, minimum), plan.max)
		plan.flex = track.Max.Fr
	case css.GridTrackMinContent:
		plan = gridContributionPlan(contribution, minimum, false)
	case css.GridTrackMaxContent:
		plan = gridContributionPlan(contribution, contribution.max, false)
	default:
		plan.max = max(contribution.max, minimum)
		plan.size = min(max(contribution.max, minimum), plan.max)
		plan.stretch = true
	}
	return plan
}
