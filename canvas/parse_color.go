package canvas

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseColor reads "#RRGGBB" or "#AARRGGBB". The hash is optional.
func ParseColor(s string) (Color, error) {
	body := strings.TrimPrefix(s, "#")
	if len(body) != 6 && len(body) != 8 {
		return 0, fmt.Errorf("antui: %q is not a colour; it looks like \"#3E63DD\"", s)
	}
	v, err := strconv.ParseUint(body, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("antui: %q is not a colour; it looks like \"#3E63DD\"", s)
	}
	if len(body) == 6 {
		v |= 0xFF000000
	}
	return Color(v), nil
}
