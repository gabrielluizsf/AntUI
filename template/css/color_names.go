package css

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// namedColors is every CSS extended colour keyword, the full CSS Color 4
// table glued to the sixteen HTML names. Each maps to the exact sRGB value a
// browser paints, so red really is the red of the rainbow and tomato really
// is a tomato. transparent and currentColor are handled as keywords by
// [ParseColor], not by this table.
var namedColors = func() map[string]canvas.Color {
	m := make(map[string]canvas.Color, len(strings.Fields(cssColorData))/2)
	fs := strings.Fields(cssColorData)
	for i := 0; i+1 < len(fs); i += 2 {
		var v uint64
		fmt.Sscanf(fs[i+1], "%x", &v)
		m[fs[i]] = canvas.Color(v | 0xFF000000)
	}
	return m
}()
