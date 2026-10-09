package css

import "github.com/gabrielluizsf/antui/canvas"

// effectStyle carries the visual-effect fields of a Style: shadows around
// and under the text, the outline that follows the box, the filters the
// painted result passes through, and the transform applied around the
// border box's TransformOrigin. OutlineStyle shares the Border* constants.
type effectStyle struct {
	BoxShadow       []Shadow
	TextShadow      []Shadow
	OutlineWidth    int
	OutlineStyle    uint8
	OutlineColor    canvas.Color
	OutlineOffset   int
	Filters         []Filter
	BackdropFilters []Filter

	// Transform is the transform list. Nil means none.
	Transform []TransformFunc
	// TransformOrigin is the transform pivot: the initial value sits in the
	// box's centre.
	TransformOrigin [2]Length
}
