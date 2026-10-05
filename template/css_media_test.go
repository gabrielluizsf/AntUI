package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// mediaPainter is a painter over an offscreen window with sheet loaded into
// its class table, the way TemplateWithCSS sets one up.
func mediaPainter(t *testing.T, win *antui.Window, sheet string) *cssStyle {
	t.Helper()
	cs := newCSSStyle(win)
	path := filepath.Join(t.TempDir(), "app.css")
	if err := os.WriteFile(path, []byte(sheet), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := cs.classes.SetStyle(path); err != nil {
		t.Fatal(err)
	}
	return cs
}

// windowBg paints the window the way a frame paints it — cleared to the
// template's background — and reads back what landed in the corner.
func windowBg(win *antui.Window, cs *cssStyle) canvas.Color {
	win.Canvas().Clear(cs.Background(win))
	return win.Canvas().At(1, 1)
}

// TestCSSMediaFollowsWindowResize is the resize the frame has to see: the
// sheet's @media windows are read against the window this frame draws at, so
// resizing between frames flips the rules, and the style cache behind the
// read is keyed by that window rather than left over from the last one.
func TestCSSMediaFollowsWindowResize(t *testing.T) {
	win := blankWin(t, 400, 400)
	cs := mediaPainter(t, win, `
		@media (max-width: 700px) { body { background-color: #0000FF; } }
		@media (min-width: 701px) { body { background-color: #FF0000; } }
	`)

	if got := windowBg(win, cs); got != canvas.RGBA(0, 0, 255, 255) {
		t.Errorf("400 wide: background %v, want the blue max-width rule", got)
	}

	win.ResizeCanvas(900, 400)
	if got := windowBg(win, cs); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("after growing to 900: background %v, want the red min-width rule", got)
	}

	win.ResizeCanvas(300, 400)
	if got := windowBg(win, cs); got != canvas.RGBA(0, 0, 255, 255) {
		t.Errorf("after shrinking to 300: background %v, want the blue rule again", got)
	}
}

// TestCSSMediaFollowsWindowHeight is the edge the width-only window could not
// see: a query over the height opens and closes as the window grows tall.
func TestCSSMediaFollowsWindowHeight(t *testing.T) {
	win := blankWin(t, 400, 400)
	cs := mediaPainter(t, win, `@media (min-height: 600px) { body { background-color: #00FF00; } }`)

	if got, want := windowBg(win, cs), win.Theme().Background; got != want {
		t.Errorf("400 tall: background %v, want the theme's %v — the rule wants 600", got, want)
	}

	win.ResizeCanvas(400, 700)
	if got := windowBg(win, cs); got != canvas.RGBA(0, 255, 0, 255) {
		t.Errorf("after growing to 700 tall: background %v, want the green min-height rule", got)
	}
}

// TestCSSMediaFollowsSystemScheme is the answer that comes from the platform
// rather than from the sheet: the window says its system paints dark, and the
// query reads that instead of treating the condition as unanswerable.
func TestCSSMediaFollowsSystemScheme(t *testing.T) {
	win := blankWin(t, 400, 400)
	cs := mediaPainter(t, win, `
		@media (prefers-color-scheme: dark) { body { background-color: #0000FF; } }
		@media not (prefers-color-scheme: dark) { body { background-color: #FF0000; } }
	`)

	// An offscreen window has no system to ask, so both sides are
	// unanswered, both rules stand, and the sheet's own order decides.
	if got := windowBg(win, cs); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("a system that has not said: background %v, want the light rule — both stand, the later one wins", got)
	}

	win.SetSystemDark(true)
	if got := windowBg(win, cs); got != canvas.RGBA(0, 0, 255, 255) {
		t.Errorf("on a dark system: background %v, want the dark rule", got)
	}

	win.SetSystemDark(false)
	if got := windowBg(win, cs); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("on a light system: background %v, want the light rule", got)
	}
}

// TestCSSMediaFollowsDisplayDensity is the resolution the frame reads off its
// display: a dense one opens the rule a 96dpi display keeps shut, and a
// display that never said leaves it standing rather than dropping it.
func TestCSSMediaFollowsDisplayDensity(t *testing.T) {
	win := blankWin(t, 400, 400)
	cs := mediaPainter(t, win, `@media (min-resolution: 144dpi) { body { background-color: #00FF00; } }`)

	if got := windowBg(win, cs); got != canvas.RGBA(0, 255, 0, 255) {
		t.Errorf("a display that has not said: background %v, want the green rule — no answer drops nothing", got)
	}

	win.SetDisplayScale(1)
	if got, want := windowBg(win, cs), win.Theme().Background; got != want {
		t.Errorf("at one pixel a point (96dpi): background %v, want the theme's %v — 144dpi is more than that", got, want)
	}

	win.SetDisplayScale(2)
	if got := windowBg(win, cs); got != canvas.RGBA(0, 255, 0, 255) {
		t.Errorf("at two pixels a point (192dpi): background %v, want the green rule", got)
	}
}
