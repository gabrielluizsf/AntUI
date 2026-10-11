package svg

// find returns the first descendant with the given tag name, or nil. It is how
// a `<defs>` or a `<style>` block is reached without walking the tree by hand.
func (e *element) find(name string) *element {
	if e.Name == name {
		return e
	}
	for _, k := range e.Kids {
		if got := k.find(name); got != nil {
			return got
		}
	}
	return nil
}
