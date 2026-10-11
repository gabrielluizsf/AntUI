package template

// multiBalance is the height each column gets when the container has no height
// of its own: the tallest column of an even split of everything the children
// asked for. A child taller than the share is not cut — there is nothing to cut
// it against — so its column is simply taller than the rest.
func (c *CSS) multiBalance(b *multiBatch) int {
	total := 0
	for i := range b.items {
		total += b.items[i].mt + b.items[i].h + b.items[i].mb
	}
	if b.count < 1 || total == 0 {
		return 0
	}
	return max((total+b.count-1)/b.count, 1)
}
