package template

// gridRowSpace is the space a row axis can fill: the declared content height
// when there is one, and nothing at all when the axis is sized by its content.
func (b *gridBatch) gridRowSpace(definite bool) int {
	if !definite {
		return 0
	}
	return b.contentH
}
