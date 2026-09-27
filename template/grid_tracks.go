package template

// Grid track sizing: what each track of a resolved axis is given, and how the
// space the container has is shared out between them.

import (
	"math"
	"sort"

	"github.com/gabrielluizsf/antui/template/css"
)

// gridTrackPlan is one track mid-resolution: the size it is heading for, the
// floor and ceiling it may not leave, and the fraction and the stretch that
// still owe it space.
type gridTrackPlan struct {
	size    int
	min     int
	max     int
	flex    float64
	stretch bool
}

// gridContribution is what the items in a track ask of it: the narrowest they
// can be and the width they would rather have.
type gridContribution struct {
	min, max int
}

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

// gridTrackMin is the floor a minmax gives its track. An auto floor is the
// automatic minimum size, which is the narrowest its items can be, and the
// content keywords are what they say they are.
func (b *gridBatch) gridTrackMin(size css.GridTrackSize, contribution gridContribution, base int, rowCyclic bool) int {
	switch size.Kind {
	case css.GridTrackLength:
		minimum := contribution.min
		if !rowCyclic || !size.Length.IsPct() {
			minimum = b.css.style.length(b.st, size.Length, base)
		}
		return max(minimum, 0)
	case css.GridTrackMinContent:
		return contribution.min
	case css.GridTrackMaxContent:
		return contribution.max
	default:
		return contribution.min
	}
}

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

// gridUnbounded is the room handed to a track that is being asked how wide it
// wants to be: no container outside has an opinion yet, so its ceiling is the
// whole int and only the track's own floor and ceiling say no.
const gridUnbounded = math.MaxInt

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

// gridTrackContributions gathers what the items spanning each track ask of it.
// An item that spans several tracks divides what it asks evenly between them,
// once for the narrowest it can be and once for what it would rather have.
func (b *gridBatch) gridTrackContributions(count int, columns bool) []gridContribution {
	contributions := make([]gridContribution, count)
	for i := range b.items {
		it := &b.items[i]
		if it.out {
			continue
		}
		start, end := it.rowStart, it.rowEnd
		minimum, size := it.contribMinH, it.contribH
		if columns {
			start, end = it.colStart, it.colEnd
			minimum, size = it.contribMinW, it.contribW
			size += it.ml + it.mr
			minimum += it.ml + it.mr
		} else {
			size += it.mt + it.mb
			minimum += it.mt + it.mb
		}
		start = max(start, 0)
		end = min(max(end, start+1), count)
		if start >= count || end <= start {
			continue
		}
		gridSpreadContribution(contributions, start, end, minimum, size)
	}
	return contributions
}

// gridSpreadContribution raises what the tracks of one item ask of them to
// cover what the item wants, the narrowest and the widest side by side.
func gridSpreadContribution(contributions []gridContribution, start, end, minimum, size int) {
	minimums := make([]int, len(contributions[start:end]))
	for i := range minimums {
		minimums[i] = contributions[start+i].min
	}
	maximums := make([]int, len(contributions[start:end]))
	for i := range maximums {
		maximums[i] = contributions[start+i].max
	}
	gridSpread(maximums, size)
	gridSpread(minimums, minimum)
	for i := range contributions[start:end] {
		contributions[start+i].min = min(minimums[i], maximums[i])
		contributions[start+i].max = maximums[i]
	}
}

// gridSpread divides what one item wants between the tracks it spans, giving
// each the same share of what is missing.
func gridSpread(values []int, value int) {
	if len(values) == 0 {
		return
	}
	used := 0
	for _, current := range values {
		used += current
	}
	if used < value {
		extra := (value - used + len(values) - 1) / len(values)
		for i := range values {
			values[i] += extra
		}
	}
}

func sumGridSizes(sizes []int) int {
	total := 0
	for _, size := range sizes {
		total += size
	}
	return total
}

func sumGridPlans(plans []gridTrackPlan) int {
	total := 0
	for _, plan := range plans {
		total += plan.size
	}
	return total
}

func sumGridPlansFor(indices []int, plans []gridTrackPlan) int {
	total := 0
	for _, index := range indices {
		total += plans[index].size
	}
	return total
}

// distributeGridFlex hands the flexible tracks the pool left over, largest
// fractional part first, then trims back down if a floor left them overshot.
func distributeGridFlex(plans []gridTrackPlan, indices []int, pool int, totalFlex float64) {
	if len(indices) == 0 || totalFlex <= 0 {
		return
	}
	fractions := make(map[int]float64, len(indices))
	for _, index := range indices {
		exact := float64(pool) * plans[index].flex / totalFlex
		size := int(math.Floor(exact))
		plans[index].size = min(max(size, plans[index].min), plans[index].max)
		fractions[index] = exact - math.Floor(exact)
	}
	remaining := pool - sumGridPlansFor(indices, plans)
	sort.SliceStable(indices, func(i, j int) bool {
		return fractions[indices[i]] > fractions[indices[j]]
	})
	for remaining > 0 {
		changed := false
		for _, index := range indices {
			if plans[index].size >= plans[index].max {
				continue
			}
			plans[index].size++
			remaining--
			changed = true
			if remaining == 0 {
				break
			}
		}
		if !changed {
			break
		}
	}
	sort.SliceStable(indices, func(i, j int) bool {
		return fractions[indices[i]] < fractions[indices[j]]
	})
	for remaining < 0 {
		changed := false
		for _, index := range indices {
			if plans[index].size <= plans[index].min {
				continue
			}
			plans[index].size--
			remaining++
			changed = true
			if remaining == 0 {
				break
			}
		}
		if !changed {
			break
		}
	}
}

// distributeGridExtra shares what is left over between the tracks that may grow
// into it, one share each and the rest to the first.
func distributeGridExtra(plans []gridTrackPlan, indices []int, extra int) {
	if len(indices) == 0 || extra <= 0 {
		return
	}
	each := extra / len(indices)
	rem := extra % len(indices)
	for i, index := range indices {
		plans[index].size += each
		if i < rem {
			plans[index].size++
		}
	}
}
