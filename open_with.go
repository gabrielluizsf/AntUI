package antui

import (
	"fmt"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// OpenWith creates and shows a window with everything Options says. Open is
// this with the defaults, and is what most games want.
//
// The window opens with the theme the system paints in, when the platform
// had an answer to give, and with [LightTheme] when it had none — after
// which [Window.Theme] is the live one and the program's to change.
func OpenWith(opt Options) (*Window, error) {
	if opt.Width < 1 || opt.Height < 1 {
		return nil, fmt.Errorf("antui: window size %dx%d is not valid", opt.Width, opt.Height)
	}
	cv, err := canvas.NewCanvas(opt.Width, opt.Height)
	if err != nil {
		return nil, err
	}
	win := &Window{
		cv:             cv,
		width:          opt.Width,
		height:         opt.Height,
		windowedWidth:  opt.Width,
		windowedHeight: opt.Height,
		limits:         opt.Limits,
		fullRedraw:     true,
		visible:        true,
		opacity:        255,
		queue:          make([]Event, 0, 32),
		theme:          LightTheme(),
	}
	win.native = newBackend()
	if err := win.native.open(win, opt.Title, opt.Width, opt.Height); err != nil {
		return nil, err
	}
	win.adoptSystemTheme()
	win.resolveGeometry(opt)

	if opt.Fullscreen {
		win.SetFullscreen(true)
	}

	win.frameStart = time.Now()
	Now()
	return win, nil
}
