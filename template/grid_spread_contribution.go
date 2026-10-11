package template

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
