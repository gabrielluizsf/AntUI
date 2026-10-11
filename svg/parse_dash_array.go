package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseDashArray reads a `stroke-dasharray`, which is a list of lengths and
// their gaps, or the word `none` for a stroke drawn whole. A list that is not a
// pattern draws no dashes rather than drawing something arbitrary, which is what
// a browser does with one it cannot read.
func parseDashArray(s string) *canvas.Dash {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "none") {
		return nil
	}
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	lengths := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, ok := parseLength(f)
		if !ok {
			return nil
		}
		lengths = append(lengths, v)
	}
	d, ok := canvas.NewDash(lengths, 0)
	if !ok {
		return nil
	}
	return &d
}
