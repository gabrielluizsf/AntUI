package antui

// Open creates and shows a window: resizable, unbounded, its size in pixels.
func Open(title string, width, height int) (*Window, error) {
	return OpenWith(Options{Title: title, Width: width, Height: height})
}
