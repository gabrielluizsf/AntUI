package template

import "github.com/gabrielluizsf/antui/template/css"

// lineHeight is the distance between two baselines: the face's own height
// enlarged by line-height, or multiplied by it when the style set one.
func (cs *cssStyle) lineHeight(st css.Style) int {
	base := cs.textOpts(st).Height()
	if st.LineHeight > 0 {
		return max(int(float64(base)*st.LineHeight+0.5), 1)
	}
	return base
}
