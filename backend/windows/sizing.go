//go:build windows

package windows

import (
	"math"
)

// sizing keeps the aspect ratio while the window is dragged. lparam points at
// the rectangle Windows is about to apply, and writing through it before
// returning is the whole mechanism.
func (d *Driver) sizing(edge uintptr, r *rect) {
	limits := d.limits
	if d.win != nil {
		limits = d.win.Bounds()
	}
	if limits.Aspect <= 0 {
		return
	}

	dx, dy := d.frame()
	width := float64(r.Right-r.Left-dx) / limits.Aspect
	height := float64(r.Bottom - r.Top - dy)

	// The axis that moves is the one the other is derived from: dragging the
	// top or the bottom sets the height and the width follows, and every
	// other edge — including a corner, where a person is thinking in width —
	// sets the width.
	switch edge {
	case wmszTop, wmszBottom:
		wanted := int32(math.Round(height*limits.Aspect)) + dx
		r.Right = r.Left + wanted
	default:
		wanted := int32(math.Round(width)) + dy
		if edge == wmszTopLeft || edge == wmszTopRight {
			r.Top = r.Bottom - wanted
		} else {
			r.Bottom = r.Top + wanted
		}
	}
}
