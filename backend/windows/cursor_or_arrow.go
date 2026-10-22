//go:build windows

package windows

// cursorOrArrow is the handle WM_SETCURSOR hands to SetCursor: the custom
// picture a caller asked for, or the ordinary system arrow once none is set.
func (d *Driver) cursorOrArrow() uintptr {
	if d.cursor != 0 {
		return d.cursor
	}
	return loadArrowCursor()
}
