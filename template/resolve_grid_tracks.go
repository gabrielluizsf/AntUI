package template

import "github.com/gabrielluizsf/antui/template/css"

// resolveGridTracks sizes every track of one axis. available is the space left
// once the gaps are taken out; definite says the axis really is that size, and
// columns picks the axis the container's content alignment stretches.
func (b *gridBatch) resolveGridTracks(tracks []css.GridTrack, columns, definite bool, available int) []int {
	if len(tracks) == 0 {
		return nil
	}
	contributions := b.gridTrackContributions(len(tracks), columns)
	plans := make([]gridTrackPlan, len(tracks))
	stretch := columns && b.st.GridJustifyContent == css.ContentStretch || !columns && b.st.AlignContent == css.ContentStretch
	base := b.contentW
	if !columns {
		base = b.contentH
		if !definite {
			base = b.contentW
		}
	}
	for i, track := range tracks {
		plans[i] = b.gridTrackPlan(track, contributions[i], base, available, definite)
	}

	nonFlex := 0
	flex := 0.0
	for _, plan := range plans {
		if plan.flex > 0 {
			flex += plan.flex
		} else {
			nonFlex += plan.size
		}
	}
	if flex > 0 && (columns || definite) {
		pool := max(available-nonFlex, 0)
		indices := make([]int, 0, len(plans))
		for i := range plans {
			if plans[i].flex > 0 {
				indices = append(indices, i)
			}
		}
		distributeGridFlex(plans, indices, pool, flex)
	} else if definite && stretch {
		used := sumGridPlans(plans)
		extra := max(available-used, 0)
		if extra > 0 {
			indices := make([]int, 0, len(plans))
			for i := range plans {
				if plans[i].stretch {
					indices = append(indices, i)
				}
			}
			distributeGridExtra(plans, indices, extra)
		}
	}

	sizes := make([]int, len(plans))
	for i := range plans {
		sizes[i] = plans[i].size
	}
	return sizes
}
