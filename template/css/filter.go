package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// Filter is one filter function. Kind and Amount carry the ones the canvas can
// apply; Drop is set for drop-shadow(), whose shape comes from the element
// rather than a matrix. A nil Drop means Kind and Amount are the filter.
type Filter struct {
	Kind   canvas.FilterKind
	Amount float64
	Drop   *Shadow
}
