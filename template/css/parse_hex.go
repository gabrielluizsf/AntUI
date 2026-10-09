package css

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// parseHex reads a hex colour: #RGB, #RGBA, #RRGGBB or #RRGGBBAA. The short
// forms double each nibble; a non-hex digit fails the whole colour.
func parseHex(s string) (canvas.Color, error) {
	body := strings.TrimPrefix(s, "#")
	valid := len(body) == 3 || len(body) == 4 || len(body) == 6 || len(body) == 8
	if !valid {
		return 0, fmt.Errorf("css: %q is not a hex colour", s)
	}
	var r, g, b, a uint8 = 0, 0, 0, 0xFF
	var err bool
	switch len(body) {
	case 3, 4:
		var ok bool
		if r, ok = hexPair(body[0], body[0]); !ok {
			err = true
		}
		if g, ok = hexPair(body[1], body[1]); !ok {
			err = true
		}
		if b, ok = hexPair(body[2], body[2]); !ok {
			err = true
		}
		if len(body) == 4 {
			if a, ok = hexPair(body[3], body[3]); !ok {
				err = true
			}
		}
	case 6, 8:
		var ok bool
		if r, ok = hexPair(body[0], body[1]); !ok {
			err = true
		}
		if g, ok = hexPair(body[2], body[3]); !ok {
			err = true
		}
		if b, ok = hexPair(body[4], body[5]); !ok {
			err = true
		}
		if len(body) == 8 {
			if a, ok = hexPair(body[6], body[7]); !ok {
				err = true
			}
		}
	}
	if err {
		return 0, fmt.Errorf("css: %q is not a hex colour", s)
	}
	return canvas.RGBA(r, g, b, a), nil
}
