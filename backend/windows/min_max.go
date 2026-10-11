//go:build windows

package windows

// minMax answers WM_GETMINMAXINFO. The numbers it wants are the whole
// window's, frame included, so the frame is added to every one of them.
func (d *Driver) minMax(mmi *minMaxInfo) {
	dx, dy := d.frame()
	limits := d.limits
	if d.win != nil {
		limits = d.win.Bounds()
	}

	// A floor under the minimum whatever was asked for: a window that cannot
	// show its own title bar is one nobody can move or close, and the
	// interface stops fitting long before that.
	minW, minH := int32(orElse(limits.MinWidth, 120)), int32(orElse(limits.MinHeight, 80))
	mmi.MinTrackSize = point{minW + dx, minH + dy}

	if limits.MaxWidth > 0 || limits.MaxHeight > 0 {
		maxW := int32(orElse(limits.MaxWidth, 1<<15))
		maxH := int32(orElse(limits.MaxHeight, 1<<15))
		mmi.MaxTrackSize = point{maxW + dx, maxH + dy}
		// MaxSize is what maximising gives, and leaving it alone would let
		// the maximise button walk straight past a limit the drag obeys.
		mmi.MaxSize = point{maxW + dx, maxH + dy}
	}
}
