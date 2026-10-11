package template

// gridInFlow is how many items the auto-placement has to find room for, which
// is what an auto-fit counts its repetitions against.
func (b *gridBatch) gridInFlow() int {
	count := 0
	for i := range b.items {
		if !b.items[i].out {
			count++
		}
	}
	return count
}
