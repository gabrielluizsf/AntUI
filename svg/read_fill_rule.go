package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// readFillRule reads `fill-rule` and `clip-rule`, which mean the same thing
// here: evenodd when written, nonzero otherwise.
func readFillRule(raw string) canvas.FillRule {
	if strings.EqualFold(strings.TrimSpace(raw), "evenodd") {
		return canvas.FillEvenOdd
	}
	return canvas.FillNonZero
}

// readLineCap reads `stroke-linecap`, defaulting to butt.
func readLineCap(raw string) canvas.LineCap {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "round":
		return canvas.CapRound
	case "square":
		return canvas.CapSquare
	}
	return canvas.CapButt
}

// readLineJoin reads `stroke-linejoin`, defaulting to miter.
func readLineJoin(raw string) canvas.LineJoin {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "round":
		return canvas.JoinRound
	case "bevel":
		return canvas.JoinBevel
	}
	return canvas.JoinMiter
}
