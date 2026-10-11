package template

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
