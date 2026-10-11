package template

import "github.com/gabrielluizsf/antui/template/css"

func gridPlacementSpan(placement css.GridPlacement) int {
	if placement.Start.Kind == css.GridLineSpan {
		return max(1, placement.Start.Index)
	}
	if placement.End.Kind == css.GridLineSpan {
		return max(1, placement.End.Index)
	}
	return 1
}
