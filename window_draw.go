package antui

import "github.com/gabrielluizsf/antui/canvas"

// Drawing on a window is drawing on its canvas. These forward so that the
// common case reads as one call rather than two.

// Clear paints the whole window in one colour.
func (win *Window) Clear(c canvas.Color) { win.cv.Clear(c) }

// Pixel blends one pixel.
func (win *Window) Pixel(x, y int, c canvas.Color) { win.cv.Pixel(x, y, c) }

// Line draws a line.
func (win *Window) Line(x0, y0, x1, y1 int, c canvas.Color) {
	win.cv.Line(x0, y0, x1, y1, c)
}

// Rect draws a one-pixel outline.
func (win *Window) Rect(x, y, w, h int, c canvas.Color) { win.cv.Rect(x, y, w, h, c) }

// FillRect fills a rectangle.
func (win *Window) FillRect(x, y, w, h int, c canvas.Color) {
	win.cv.FillRect(x, y, w, h, c)
}

// RoundRect draws a rounded outline.
func (win *Window) RoundRect(x, y, w, h, r int, c canvas.Color) {
	win.cv.RoundRect(x, y, w, h, r, c)
}

// FillRoundRect fills a rounded rectangle.
func (win *Window) FillRoundRect(x, y, w, h, r int, c canvas.Color) {
	win.cv.FillRoundRect(x, y, w, h, r, c)
}

// Circle draws a circle outline.
func (win *Window) Circle(cx, cy, radius int, c canvas.Color) {
	win.cv.Circle(cx, cy, radius, c)
}

// FillCircle fills a disc.
func (win *Window) FillCircle(cx, cy, radius int, c canvas.Color) {
	win.cv.FillCircle(cx, cy, radius, c)
}

// Triangle draws the edges of a triangle.
func (win *Window) Triangle(x0, y0, x1, y1, x2, y2 int, c canvas.Color) {
	win.cv.Triangle(x0, y0, x1, y1, x2, y2, c)
}

// FillTriangle fills a triangle.
func (win *Window) FillTriangle(x0, y0, x1, y1, x2, y2 int, c canvas.Color) {
	win.cv.FillTriangle(x0, y0, x1, y1, x2, y2, c)
}

// Text draws a string and returns the width of the widest line drawn.
func (win *Window) Text(x, y int, text string, c canvas.Color) int {
	return win.cv.Text(x, y, text, c)
}

// TextScaled draws a string with the font blown up, returning the width of
// the widest line drawn.
func (win *Window) TextScaled(x, y int, text string, c canvas.Color, scale int) int {
	return win.cv.TextScaled(x, y, text, c, scale)
}

// Blit composites another canvas over the window.
func (win *Window) Blit(x, y int, src *canvas.Canvas) { win.cv.Blit(x, y, src) }

// BlitScaled draws another canvas fitted into a rectangle.
func (win *Window) BlitScaled(x, y, width, height int, src *canvas.Canvas) {
	win.cv.BlitScaled(x, y, width, height, src)
}

// SetClip narrows drawing to a rectangle.
func (win *Window) SetClip(x, y, w, h int) { win.cv.SetClip(x, y, w, h) }

// ResetClip opens drawing back up to the whole window.
func (win *Window) ResetClip() { win.cv.ResetClip() }