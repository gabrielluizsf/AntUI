package backend

// Tool is what is touching the screen.
type Tool int

// The tools. A device that does not say reports ToolUnknown rather than
// guessing at a finger.
const (
	ToolUnknown Tool = iota
	ToolFinger
	ToolStylus
	ToolMouse
	ToolEraser // the far end of a stylus, used to rub out
)

func (t Tool) String() string {
	switch t {
	case ToolFinger:
		return "finger"
	case ToolStylus:
		return "stylus"
	case ToolMouse:
		return "mouse"
	case ToolEraser:
		return "eraser"
	}
	return "unknown"
}