package template

import "github.com/gabrielluizsf/antui/template/css"

// gridTrackPlan sizes one track. A track written as a length is exactly that
// length; a fraction grows into what the flexible tracks share; auto takes what
// its items ask for and stretches into the space left over; and the content
// keywords ask for the narrowest or the widest they can be, with fit-content()
// clamped between the two by the limit it was given.
func (b *gridBatch) gridTrackPlan(track css.GridTrack, contribution gridContribution, base, available int, definite bool) gridTrackPlan {
	rowCyclic := !definite
	if track.Kind == css.GridTrackSingle {
		switch track.Size.Kind {
		case css.GridTrackLength:
			size := contribution.max
			if !rowCyclic || !track.Size.Length.IsPct() {
				size = b.css.style.length(b.st, track.Size.Length, base)
			}
			return gridTrackPlan{size: max(size, 0), min: max(size, 0), max: max(size, 0)}
		case css.GridTrackFlexible:
			maxSize := contribution.max
			if definite {
				maxSize = max(available, contribution.max)
			}
			return gridTrackPlan{size: contribution.max, min: 0, max: maxSize, flex: track.Size.Fr}
		case css.GridTrackMinContent:
			return gridContributionPlan(contribution, contribution.min, false)
		case css.GridTrackMaxContent:
			return gridContributionPlan(contribution, contribution.max, false)
		case css.GridTrackFitContent:
			limit := max(b.css.style.length(b.st, track.Size.Length, base), 0)
			return gridContributionPlan(contribution, clampGridItem(contribution.max, contribution.min, limit), false)
		default:
			return gridTrackPlan{size: contribution.max, min: contribution.min, max: contribution.max, stretch: true}
		}
	}

	minimum := b.gridTrackMin(track.Min, contribution, base, rowCyclic)
	plan := gridTrackPlan{min: minimum}
	switch track.Max.Kind {
	case css.GridTrackLength:
		maximum := contribution.max
		if !rowCyclic || !track.Max.Length.IsPct() {
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
