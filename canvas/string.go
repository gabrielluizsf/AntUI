package canvas

import "fmt"

// String writes a colour the way it is written everywhere else: "#RRGGBB",
// or "#AARRGGBB" when it is not opaque.
func (c Color) String() string {
	if c.A() == 0xFF {
		return fmt.Sprintf("#%02X%02X%02X", c.R(), c.G(), c.B())
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.A(), c.R(), c.G(), c.B())
}
