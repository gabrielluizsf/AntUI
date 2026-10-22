//go:build windows

package windows

// idcArrow is the system's ordinary pointer, picked by name from the
// resources built into Windows itself rather than drawn by this package.
const idcArrow = 32512

// loadArrowCursor is the cursor a window gets back once a custom one is let
// go of. It belongs to Windows, not this process, so there is nothing to
// destroy when done with it.
func loadArrowCursor() uintptr {
	cursor, _, _ := procLoadCursorW.Call(0, idcArrow)
	return cursor
}
