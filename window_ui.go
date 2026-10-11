package antui

// windowUI is the immediate-mode widget state: hover, press, focus and the
// live text edit.
type windowUI struct {
	uiHot    uint32 // the widget under the pointer
	uiActive uint32 // the widget being pressed
	uiFocus  uint32 // the widget receiving the keyboard
	uiCursor int    // cursor position in the focused text field
	uiBlink  float64
	// editBuf splices typed text into a field; the window keeps it so a
	// keystroke does not allocate a builder every frame.
	editBuf []byte

	// uiTab is the widgets that take the keyboard, in the order they drew
	// this frame; Tab walks it. Its slice is reused from frame to frame.
	uiTab []uint32
}
