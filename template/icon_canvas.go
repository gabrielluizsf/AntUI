package template

import "github.com/gabrielluizsf/antui/canvas"

// IconCanvas is the picture a named icon makes at a size in a colour, already
// painted and ready to be blitted, or nil when there is no such icon. It is
// what [DrawIcon] draws and what a stylesheet's `url()` puts on a box, so both
// ways of showing an icon share one picture and one cache.
//
// The canvas is the icon itself and belongs to the registry: blit from it as
// many times as a frame needs and do not draw on it.
func IconCanvas(name string, area canvas.Area, tint canvas.Color) *canvas.Canvas {
	w, h := area.Width, area.Height
	if w <= 0 || h <= 0 {
		return nil
	}
	iconMu.RLock()
	e := icons[name]
	iconMu.RUnlock()
	if e == nil || e.img == nil {
		return nil
	}
	// A drawing that moves with the clock has no picture worth remembering
	// here: the canvas of the last frame is a frame that has ended. It goes
	// back to the drawing every time, and the drawing holds the moment — this
	// only asks for it.
	if e.img.Animates() {
		return e.img.RenderWith(w, h, tint)
	}
	key := iconPaintKey{name: name, w: w, h: h, tint: tint}
	iconPaintMu.RLock()
	painted, ok := iconPaint[key]
	iconPaintMu.RUnlock()
	if ok {
		return painted
	}
	painted = e.img.RenderWith(w, h, tint)
	if painted == nil {
		return nil
	}
	iconPaintMu.Lock()
	// Starting over rather than keeping the oldest: an interface asks for the
	// same few pictures every frame, so what is worth keeping is decided by this
	// frame and not by what happened to be asked for first.
	if len(iconPaint) >= iconPaintLimit {
		iconPaint = make(map[iconPaintKey]*canvas.Canvas, iconPaintLimit)
	}
	iconPaint[key] = painted
	iconPaintMu.Unlock()
	return painted
}
