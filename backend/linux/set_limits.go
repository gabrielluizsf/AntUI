//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"github.com/gabrielluizsf/antui/backend"
)

// SetLimits passes the size constraints on to the window manager as
// WM_NORMAL_HINTS. A window with no bounds at all still gets the property
// written, because leaving a stale one behind would lock a window that had
// just been unlocked.
func (x *Driver) SetLimits(limits backend.Limits) {
	if !x.alive || x.window == 0 {
		return
	}

	// Eighteen 32-bit fields: flags, four obsolete position and size fields
	// no one has read since 1989, then the sizes, the increments, the two
	// aspect ratios, the base size and the gravity.
	hints := make([]uint32, 18)
	if limits.MinWidth > 0 || limits.MinHeight > 0 {
		hints[0] |= xSizeHintMinSize
		hints[5] = uint32(max(limits.MinWidth, 1))
		hints[6] = uint32(max(limits.MinHeight, 1))
	}
	if limits.MaxWidth > 0 || limits.MaxHeight > 0 {
		hints[0] |= xSizeHintMaxSize
		// A maximum on one axis only still has to name both, so the axis
		// with no maximum is given one no display will reach.
		hints[7] = uint32(orElse(limits.MaxWidth, 1<<15))
		hints[8] = uint32(orElse(limits.MaxHeight, 1<<15))
	}
	if limits.Aspect > 0 {
		// The ratio goes as a fraction of two integers, and both ends are
		// set to it: the window manager keeps the shape between them, so an
		// equal pair is how an exact ratio is said.
		num, den := ratio(limits.Aspect)
		hints[0] |= xSizeHintAspect
		hints[11], hints[12] = num, den // smallest allowed
		hints[13], hints[14] = num, den // largest allowed
	}
	if limits.MinWidth > 0 && limits.MinHeight > 0 {
		// Without a base size the increments and the aspect are measured
		// from zero, which is not what a minimum meant.
		hints[0] |= xSizeHintBaseSize
		hints[15] = uint32(limits.MinWidth)
		hints[16] = uint32(limits.MinHeight)
	}

	x.changeProperty(xWMNormalHints, xWMSizeHints, 32, le32(hints...), 18)
}
