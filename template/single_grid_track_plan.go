package template

import "github.com/gabrielluizsf/antui/template/css"

// singleGridTrackPlan sizes a track written as one value — a length, a
// fraction, or a content keyword — rather than as a min/max pair.
func (b *gridBatch) singleGridTrackPlan(track css.GridTrack, contribution gridContribution, base, available int, definite bool) gridTrackPlan {
	switch track.Size.Kind {
	case css.GridTrackLength:
		size := contribution.max
		if definite || !track.Size.Length.IsPct() {
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
