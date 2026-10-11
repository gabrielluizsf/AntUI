package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// loadingPainter is the half of a template the loading screen needs: the
// colour the window is cleared to, and the progress bar it draws with. Both
// [Template] and *uiTemplate answer to it, so the screen is drawn in the
// voice of whichever template the program picked — the same bar the rest of
// the app wears.
type loadingPainter interface {
	Background(win *antui.Window) canvas.Color
	Progress(win *antui.Window, x, y, w, h int, progress float32)
}
