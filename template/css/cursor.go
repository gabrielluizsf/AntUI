package css

// Cursor names the pointer shape a box asks for. The engine only models it;
// painting a real OS cursor is the window's job.
const (
	CursorDefault uint8 = iota
	CursorPointer
	CursorText
	CursorWait
	CursorCrosshair
	CursorMove
	CursorNotAllowed
	CursorGrab
	CursorGrabbing
	CursorCell
)

func parseCursor(raw string) (uint8, bool) {
	switch raw {
	case "default", "auto":
		return CursorDefault, true
	case "pointer":
		return CursorPointer, true
	case "text":
		return CursorText, true
	case "wait", "progress":
		return CursorWait, true
	case "crosshair":
		return CursorCrosshair, true
	case "move", "all-scroll":
		return CursorMove, true
	case "not-allowed":
		return CursorNotAllowed, true
	case "grab":
		return CursorGrab, true
	case "grabbing":
		return CursorGrabbing, true
	case "cell":
		return CursorCell, true
	}
	return 0, false
}
