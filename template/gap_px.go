package template

import "github.com/gabrielluizsf/antui/template/css"

// gapPx resolves one gap length. Percentages measure the container's content
// width; CSS ties row-gap to the content height, but the height only exists
// once solving finishes, so the width stands in for both.
func (b *flexBatch) gapPx(l css.Length) int {
	return b.css.style.length(b.st, l, b.contentW)
}
