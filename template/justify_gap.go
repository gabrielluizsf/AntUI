package template

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// justifyGap is the extra space each gap takes when a line is justified to
// fill avail, or zero when it cannot grow.
func justifyGap(cv *canvas.Canvas, text string, o canvas.TextStyle, avail int) int {
	gaps := strings.Count(text, " ")
	if gaps == 0 {
		return 0
	}
	if extra := (avail - cv.TextWidthStyled(text, o)) / gaps; extra > 0 {
		return extra
	}
	return 0
}
