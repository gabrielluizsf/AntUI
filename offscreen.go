package antui

import (
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// Offscreen makes a Window that draws into a canvas and shows nothing: no
// window is opened, nothing is presented, and events never come from a
// display. It is for tests, and for rendering a UI onto a texture to be
// pasted somewhere else. `win.Canvas()` is where everything lands.
//
// Its input still works — [Window.Push] folds events into the same input
// state a real window keeps — so a widget can be clicked and typed at from a
// test exactly as from a screen. [Window.End] does nothing without a display,
// and the frame is free to be looked at after [Window.Begin].
func Offscreen(width, height int) (*Window, *canvas.Canvas, error) {
	cv, err := canvas.NewCanvas(width, height)
	if err != nil {
		return nil, nil, err
	}
	win := &Window{
		cv:             cv,
		width:          width,
		height:         height,
		windowedWidth:  width,
		windowedHeight: height,
		fullRedraw:     true,
		visible:        true,
		opacity:        255,
		queue:          make([]Event, 0, 32),
		theme:          LightTheme(),
		frameStart:     time.Now(),
	}
	Now()
	return win, cv, nil
}
