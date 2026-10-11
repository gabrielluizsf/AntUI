package template

import (
	"math"
	"sort"
)

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
