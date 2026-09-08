package antui

import (
	"fmt"
	"iter"
	"strings"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// backend is what each platform implements. The window core above it knows
// nothing about X11, Win32 or Cocoa; it hands the backend a finished frame
// and asks it for events.
//
// It lives in this package and not in a backend package of its own because it
// is both handed a *Window to pump events into and called "on" that Window —
// the two halves are mutually recursive by design. A platform's work lives in
// backend/linux, backend/macos and backend/windows; the
// thin backend_*.go files in this package sit between them and this
// interface, turning each platform package's exported calls into what Window
// above expects.
type platform interface {
	open(win *Window, title string, width, height int) error
	close()
	pump(win *Window)
	present(win *Window, dirty canvas.Area)
	setTitle(title string)
	// setFullscreen reports whether the request could be made at all, not
	// whether the window manager honoured it.
	setFullscreen(on bool) bool
	displaySize() (w, h int, ok bool)
	displayRefresh() int
	// setLimits passes the size constraints on to the window manager. They
	// are advice on X11 and enforced on Win32 and Cocoa; see window_size.go.
	setLimits(win *Window, limits Limits)
	// setSize asks for a new drawable size, reporting whether the request
	// could be made rather than whether it was honoured.
	setSize(win *Window, width, height int) bool
	// contentScale is pixels per point, or 0 when the system does not say.
	contentScale() float64
	// clipboard is what the system clipboard holds: its text, and the files
	// it names when it names any. Both empty when there is nothing there or
	// the platform cannot be asked.
	clipboard() (text string, files []string)
	// setIcon gives the window an icon, largest first. It reports whether the
	// window system took it.
	setIcon(images []*canvas.Canvas) bool
	// setClipboard puts text on the system clipboard, reporting whether the
	// system took it.
	setClipboard(text string) bool
}

// textBuffer bounds the text typed in one frame, in UTF-8 bytes.
const textBuffer = 64

// Window is an open window and everything drawn into it. It is not safe for
// concurrent use: like every immediate-mode UI, one goroutine owns the frame
// loop and everything else talks to it through channels.
type Window struct {
	cv     *canvas.Canvas
	native platform

	shadow                    []canvas.Color // a copy of the previous frame
	shadowWidth, shadowHeight int

	width, height int
	shouldClose   bool
	fullRedraw    bool // forces sending the whole frame
	visible       bool

	fullscreen                    bool // what was last asked for
	windowedWidth, windowedHeight int  // the size to come back out to

	limits Limits  // what the window may be resized to; see window_size.go
	scale  float64 // pixels per point, 0 until the backend says

	queue         []Event
	droppedEvents int

	mouseX, mouseY   int
	mouseDX, mouseDY int
	mouseState       [mouseCount]bool
	mouseDownFrame   [mouseCount]bool
	mouseUpFrame     [mouseCount]bool
	wheel            int
	mods             Mod

	keyState     [keyCount]bool
	keyDownFrame [keyCount]bool
	keyUpFrame   [keyCount]bool

	text strings.Builder
	// dropped is what was dragged onto the window this frame.
	dropped []string
	// touchFirst is a platform with nothing but a touchscreen, and sawTouch
	// is a finger having actually arrived. See Window.TouchScreen.
	touchFirst bool
	sawTouch   bool
	// touches is every finger on the screen; see touch.go.
	touches []Touch
	// gesture is what those fingers are in the middle of doing; see
	// gesture.go.
	gesture gestureState
	// safe is the part of the canvas the system is not covering; see
	// safearea.go. Zero means all of it.
	safe canvas.Area

	frameStart time.Time
	delta      float64
	targetFPS  int

	theme    Theme
	uiHot    uint32 // the widget under the pointer
	uiActive uint32 // the widget being pressed
	uiFocus  uint32 // the widget receiving the keyboard
	uiCursor int    // cursor position in the focused text field
	uiBlink  float64
}

var (
	clockOnce sync.Once
	clockZero time.Time
)

// Now is the seconds since the library was first asked the time, on a
// monotonic clock that no change to the system clock can move.
func Now() float64 {
	clockOnce.Do(func() { clockZero = time.Now() })
	return time.Since(clockZero).Seconds()
}

// Sleep pauses for a number of seconds.
func Sleep(seconds float64) {
	if seconds > 0 {
		time.Sleep(time.Duration(seconds * float64(time.Second)))
	}
}

// OpenWith creates and shows a window with everything Options says. Open is
// this with the defaults, and is what most games want.
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
		queue:          make([]Event, 0, 32),
		theme:          LightTheme(),
	}
	win.native = newBackend()
	if err := win.native.open(win, opt.Title, opt.Width, opt.Height); err != nil {
		return nil, err
	}

	win.scale = win.native.contentScale()
	width, height := opt.Width, opt.Height
	if opt.Points && win.scale > 0 {
		width, height = scaleTo(width, win.scale), scaleTo(height, win.scale)
		win.limits = win.limits.Scaled(win.scale)
	}
	if dw, dh, ok := win.native.displaySize(); ok {
		width, height = fitDisplay(width, height, dw, dh)
	}
	width, height = win.limits.Clamp(width, height)

	if width != opt.Width || height != opt.Height {
		win.native.setSize(win, width, height)
		win.ResizeCanvas(width, height)
		win.windowedWidth, win.windowedHeight = width, height
	}
	win.native.setLimits(win, win.Bounds())
	if opt.Fullscreen {
		win.SetFullscreen(true)
	}

	win.frameStart = time.Now()
	Now()
	return win, nil
}

// Close destroys the window and gives back everything it held.
func (win *Window) Close() {
	if win == nil || win.native == nil {
		return
	}
	win.native.close()
	win.native = nil
}

// Begin starts a frame: it processes system events and refreshes the input
// state. It returns false once the window has been closed.
//
//	for win.Begin() {
//	    win.Clear(canvas.White)
//	    win.End()
//	}
func (win *Window) Begin() bool {
	now := time.Now()
	win.delta = now.Sub(win.frameStart).Seconds()
	if win.delta < 0 || win.delta > 1 {
		win.delta = 1.0 / 60.0
	}
	win.frameStart = now

	win.mouseDownFrame = [mouseCount]bool{}
	win.mouseUpFrame = [mouseCount]bool{}
	win.keyDownFrame = [keyCount]bool{}
	win.keyUpFrame = [keyCount]bool{}
	win.wheel = 0
	win.mouseDX, win.mouseDY = 0, 0
	win.text.Reset()
	win.dropped = win.dropped[:0]
	win.beginTouchFrame()
	win.queue = win.queue[:0]
	win.droppedEvents = 0
	win.uiHot = 0
	win.uiBlink += win.delta

	if win.native != nil {
		win.native.pump(win)
	}
	win.recognize(Now())

	win.cv.ResetClip()
	return !win.shouldClose
}

// End finishes the frame: it sends the pixels to the screen and applies the
// frame-rate cap.
func (win *Window) End() {
	if win.native == nil {
		return
	}
	win.native.present(win, win.dirtyRegion())

	if win.targetFPS > 0 {
		budget := time.Duration(float64(time.Second) / float64(win.targetFPS))
		if elapsed := time.Since(win.frameStart); elapsed < budget {
			time.Sleep(budget - elapsed)
		}
	}
}

// Quit marks the window for closing, so the next Begin returns false.
func (win *Window) Quit() { win.shouldClose = true }

// Running reports whether the window is still open.
func (win *Window) Running() bool { return win != nil && !win.shouldClose }

// SetTitle changes the title bar.
func (win *Window) SetTitle(title string) {
	if win.native != nil {
		win.native.setTitle(title)
	}
}

// Width is the drawable width in pixels.
func (win *Window) Width() int { return win.cv.Width }

// Height is the drawable height in pixels.
func (win *Window) Height() int { return win.cv.Height }

// Delta is how long the last frame took, in seconds.
func (win *Window) Delta() float64 { return win.delta }

// SetFPS caps the frame rate in End. Zero turns the cap off.
func (win *Window) SetFPS(fps int) { win.targetFPS = max(fps, 0) }

// SetFullscreen asks the system to put the window full screen, or to take it
// back out. It reports whether the request could be made, not whether it was
// granted.
func (win *Window) SetFullscreen(on bool) bool {
	if on && !win.fullscreen {
		win.windowedWidth = win.cv.Width
		win.windowedHeight = win.cv.Height
	}
	if win.native == nil || !win.native.setFullscreen(on) {
		return false
	}
	win.fullscreen = on
	win.fullRedraw = true
	return true
}

// Fullscreen reports what was last asked for, not what a window manager did
// about it.
func (win *Window) Fullscreen() bool { return win.fullscreen }

// DisplaySize is the display's size in pixels, and whether it could be told.
func (win *Window) DisplaySize() (width, height int, ok bool) {
	if win.native == nil {
		return 0, 0, false
	}
	return win.native.displaySize()
}

// DisplayRefresh is how many times a second the display refreshes, or 0 when
// it could not be told.
func (win *Window) DisplayRefresh() int {
	if win.native == nil {
		return 0
	}
	return win.native.displayRefresh()
}

// Canvas is the window's pixel surface: direct access, if you need it.
func (win *Window) Canvas() *canvas.Canvas { return win.cv }

// Theme returns the live theme. Change its fields to restyle the widgets.
func (win *Window) Theme() *Theme { return &win.theme }

// SetTheme replaces the whole theme at once.
func (win *Window) SetTheme(theme Theme) { win.theme = theme }

// ---------------------------------------------------------------------------
// Events
// ---------------------------------------------------------------------------

// push folds an event into the input state and queues it for NextEvent.
func (win *Window) Push(ev Event) {
	switch ev.Type {
	case EventClose:
		win.shouldClose = true
	case EventMouseMove:
		win.mouseDX += ev.X - win.mouseX
		win.mouseDY += ev.Y - win.mouseY
		win.mouseX, win.mouseY = ev.X, ev.Y
	case EventMouseDown:
		win.mouseX, win.mouseY = ev.X, ev.Y
		if ev.Button >= 0 && ev.Button < mouseCount {
			win.mouseState[ev.Button] = true
			win.mouseDownFrame[ev.Button] = true
		}
		win.mods = ev.Mods
	case EventMouseUp:
		win.mouseX, win.mouseY = ev.X, ev.Y
		if ev.Button >= 0 && ev.Button < mouseCount {
			win.mouseState[ev.Button] = false
			win.mouseUpFrame[ev.Button] = true
		}
		win.mods = ev.Mods
	case EventMouseWheel:
		win.wheel += ev.Wheel
	case EventKeyDown:
		if ev.Key > 0 && ev.Key < keyCount {
			win.keyState[ev.Key] = true
			win.keyDownFrame[ev.Key] = true
		}
		win.mods = ev.Mods
	case EventKeyUp:
		if ev.Key > 0 && ev.Key < keyCount {
			win.keyState[ev.Key] = false
			win.keyUpFrame[ev.Key] = true
		}
		win.mods = ev.Mods
	case EventText:
		if win.text.Len()+len(ev.Text) < textBuffer {
			win.text.WriteString(ev.Text)
		}
	case EventDropFiles:
		win.mouseX, win.mouseY = ev.X, ev.Y
		win.dropped = append(win.dropped, ev.Files...)
	case EventExpose:
		win.fullRedraw = true
	case EventFocus:
		if !ev.Focused {
			win.keyState = [keyCount]bool{}
			win.mouseState = [mouseCount]bool{}
			win.mods = 0
			for i := range win.touches {
				win.touches[i].Ended = true
				win.touches[i].Cancelled = true
			}
		}
	case EventTouchDown, EventTouchMove, EventTouchUp, EventTouchCancel:
		win.pushTouch(ev)
	}

	if len(win.queue) >= maxEvents {
		win.droppedEvents++
		return
	}
	win.queue = append(win.queue, ev)
}

// pushSimple queues an event that carries nothing but its kind.
func (win *Window) PushSimple(t EventType) {
	win.Push(Event{Type: t, Mods: win.mods, X: win.mouseX, Y: win.mouseY})
}

// The methods below are the half of Window a backend is allowed to touch,
// through backend.Face. They are exported, because an interface with an
// unexported method can only be implemented in the package that defines it
// and a platform backend lives in its own folder — but a program has no
// reason to call them; a backend's Driver receives the Window through the
// interface and drives it with these.

// SetMouse is where a backend with no cursor position of its own records the
// pointer.
func (win *Window) SetMouse(x, y int) { win.mouseX, win.mouseY = x, y }

// SetShouldClose asks the window to begin closing.
func (win *Window) SetShouldClose() { win.shouldClose = true }

// NextEvent pops the next event off this frame's queue, reporting false when
// it is empty.
func (win *Window) NextEvent() (Event, bool) {
	if len(win.queue) == 0 {
		return Event{}, false
	}
	ev := win.queue[0]
	win.queue = win.queue[1:]
	return ev, true
}

// Events iterates this frame's queue, draining it as it goes.
func (win *Window) Events(yield func(Event) bool) {
	for {
		ev, ok := win.NextEvent()
		if !ok || !yield(ev) {
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Input
// ---------------------------------------------------------------------------

// MouseX is the pointer's horizontal position, in window pixels.
func (win *Window) MouseX() int { return win.mouseX }

// MouseY is the pointer's vertical position, in window pixels.
func (win *Window) MouseY() int { return win.mouseY }

// MouseDown reports whether a button is held.
func (win *Window) MouseDown(b MouseButton) bool {
	return b >= 0 && b < mouseCount && win.mouseState[b]
}

// MousePressed reports whether a button went down during this frame.
func (win *Window) MousePressed(b MouseButton) bool {
	return b >= 0 && b < mouseCount && win.mouseDownFrame[b]
}

// MouseReleased reports whether a button came up during this frame.
func (win *Window) MouseReleased(b MouseButton) bool {
	return b >= 0 && b < mouseCount && win.mouseUpFrame[b]
}

// Wheel is the wheel movement this frame, positive upwards.
func (win *Window) Wheel() int { return win.wheel }

// Keys walks every key that is held or that went down or up this frame.
func (win *Window) Keys() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		if win == nil {
			return
		}
		for k := range keyCount {
			if !win.keyState[k] && !win.keyDownFrame[k] && !win.keyUpFrame[k] {
				continue
			}
			if !yield(Key(k)) {
				return
			}
		}
	}
}

// DroppedFiles is what was dragged onto the window this frame.
func (win *Window) DroppedFiles() []string {
	if win == nil {
		return nil
	}
	return win.dropped
}

// ClipboardText is what the system clipboard holds as text, or "" when it
// holds nothing this window can read.
func (win *Window) ClipboardText() string {
	if win == nil || win.native == nil {
		return ""
	}
	text, _ := win.native.clipboard()
	return text
}

// SetClipboardText puts text on the system clipboard — what a Copy button
// does — and reports whether the system took it.
func (win *Window) SetClipboardText(text string) bool {
	if win == nil || win.native == nil {
		return false
	}
	return win.native.setClipboard(text)
}

// ClipboardFiles is the files the clipboard names. Empty when it holds none.
func (win *Window) ClipboardFiles() []string {
	if win == nil || win.native == nil {
		return nil
	}
	_, files := win.native.clipboard()
	return files
}

// KeyDown reports whether a key is held.
func (win *Window) KeyDown(k Key) bool {
	return k > 0 && k < keyCount && win.keyState[k]
}

// KeyPressed reports whether a key went down during this frame.
func (win *Window) KeyPressed(k Key) bool {
	return k > 0 && k < keyCount && win.keyDownFrame[k]
}

// KeyReleased reports whether a key came up during this frame.
func (win *Window) KeyReleased(k Key) bool {
	return k > 0 && k < keyCount && win.keyUpFrame[k]
}

// Mods is the set of modifier keys held.
func (win *Window) Mods() Mod { return win.mods }

// TextInput is the text typed this frame, in UTF-8, empty when nothing was.
func (win *Window) TextInput() string { return win.text.String() }

// ---------------------------------------------------------------------------
// Resizing and the dirty region
// ---------------------------------------------------------------------------

// resizeCanvas grows or shrinks the drawable surface. Backends call it when
// the window system reports a new size.
func (win *Window) ResizeCanvas(width, height int) bool {
	width, height = max(width, 1), max(height, 1)
	if width == win.cv.Width && height == win.cv.Height {
		return true
	}
	cv, err := canvas.NewCanvas(width, height)
	if err != nil {
		return false
	}
	win.cv = cv
	win.width, win.height = width, height

	win.shadow = nil
	win.shadowWidth, win.shadowHeight = 0, 0
	win.fullRedraw = true
	win.safe = canvas.Area{}
	return true
}

// dirtyRegion is the rectangle that changed since the previous frame, so that
// only that much is sent to the display.
func (win *Window) dirtyRegion() canvas.Area {
	cv := win.cv
	all := canvas.Area{X: 0, Y: 0, Width: cv.Width, Height: cv.Height}
	count := cv.Width * cv.Height

	if win.fullRedraw || win.shadow == nil ||
		win.shadowWidth != cv.Width || win.shadowHeight != cv.Height {
		win.shadow = make([]canvas.Color, count)
		copy(win.shadow, cv.Pixels)
		win.shadowWidth, win.shadowHeight = cv.Width, cv.Height
		win.fullRedraw = false
		return all
	}

	firstRow, lastRow := -1, -1
	for y := range cv.Height {
		a := cv.Pixels[y*cv.Stride : y*cv.Stride+cv.Width]
		b := win.shadow[y*cv.Width : y*cv.Width+cv.Width]
		if !equalRow(a, b) {
			if firstRow < 0 {
				firstRow = y
			}
			lastRow = y
		}
	}
	if firstRow < 0 {
		return canvas.Area{}
	}

	minX, maxX := cv.Width, -1
	for y := firstRow; y <= lastRow; y++ {
		a := cv.Pixels[y*cv.Stride : y*cv.Stride+cv.Width]
		b := win.shadow[y*cv.Width : y*cv.Width+cv.Width]
		for x := 0; x < minX; x++ {
			if a[x] != b[x] {
				minX = x
				break
			}
		}
		for x := cv.Width - 1; x > maxX; x-- {
			if a[x] != b[x] {
				maxX = x
				break
			}
		}
		if minX == 0 && maxX == cv.Width-1 {
			break
		}
	}
	if maxX < minX {
		return canvas.Area{}
	}

	for y := firstRow; y <= lastRow; y++ {
		copy(win.shadow[y*cv.Width:y*cv.Width+cv.Width],
			cv.Pixels[y*cv.Stride:y*cv.Stride+cv.Width])
	}
	return canvas.Area{
		X: minX, Y: firstRow,
		Width: maxX - minX + 1, Height: lastRow - firstRow + 1,
	}
}

func equalRow(a, b []canvas.Color) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}