//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// ContentScale is how much the display is scaled by, or 0 when it cannot be
// known: [Driver.DisplayScale] with "the display did not say" left as the
// zero that arithmetic wants rather than as a bool to carry.
func (x *Driver) ContentScale() float64 {
	scale, ok := x.DisplayScale()
	if !ok {
		return 0
	}
	return scale
}
