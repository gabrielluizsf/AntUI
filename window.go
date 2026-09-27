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
	// setOpacity fades the whole window — decorations and all — toward
	// transparent, so what sits behind the program shows through. 255 is
	// fully opaque. It reports whether the request could be made, not whether
	// a compositor honoured it.
	setOpacity(alpha uint8) bool
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

// framePacer is the optional half of a backend, and it is told when the next
// frame is due rather than how long a frame may take: a platform that has to
// wait for something to happen can spend the time until then asleep in the
// kernel instead of coming straight back to be called again — the difference
// between an idle window costing nothing and an idle window costing a core.
// One that cannot wait is unaffected, since the core does the waiting itself.
type framePacer interface {
	setFrameDeadline(due time.Time)
}

// textBuffer bounds the text typed in one frame, in UTF-8 bytes.
const textBuffer = 64

// fallbackFPS is the frame rate a window keeps when the display will not say
// how often it refreshes: a server without RandR, an offscreen canvas, a
// platform still to be written.
const fallbackFPS = 60

// Window is an open window and everything drawn into it. It is not safe for
// concurrent use: like every immediate-mode UI, one goroutine owns the frame
// loop and everything else talks to it through channels.
type Window struct {
	cv     *canvas.Canvas
	native platform

	// shadow is the frame that was last presented, which the canvas checks
	// every write against: what the screen already holds is what does not
	// have to be sent again. It is grown once and kept between resizes.
	shadow []canvas.Color

	width, height int
	shouldClose   bool
	fullRedraw    bool // forces sending the whole frame
	presented     bool // whether the last End reached the display
	visible       bool

	fullscreen                    bool  // what was last asked for
	windowedWidth, windowedHeight int   // the size to come back out to
	opacity                       uint8 // what was last asked for; 255 = opaque

	limits Limits  // what the window may be resized to; see window_size.go
	scale  float64 // pixels per point, 0 until the backend says

	queue         []Event
	droppedEvents int

	// frontmost are draws that run at the very end of the frame, after every
	// widget painted in flow order, so an open dropdown or calendar stays on
	// top of whatever the layout drew under it.
	frontmost []func()

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
	nextFrame  time.Time // when the next frame is due; see pace
	targetFPS  int       // 0 follows the display, see SetFPS
	refreshHz  int       // the display's rate, 0 until asked, -1 when unknown
	sleep      func(time.Duration)

	theme    Theme
	uiHot    uint32 // the widget under the pointer
	uiActive uint32 // the widget being pressed
	uiFocus  uint32 // the widget receiving the keyboard
	uiCursor int    // cursor position in the focused text field
	uiBlink  float64

	// uiTab is the widgets that take the keyboard, in the order they drew
	// this frame; Tab walks it. Its slice is reused from frame to frame.
	uiTab []uint32
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
		opacity:        255,
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
	win.uiTab = win.uiTab[:0]
	win.uiBlink += win.delta
	win.frontmost = win.frontmost[:0]

	if win.native != nil {
		win.native.pump(win)
	}
	win.recognize(Now())

	win.cv.ResetClip()
	return !win.shouldClose
}

// End finishes the frame: it sends the pixels that changed to the screen, then
// holds the loop to the frame rate.
//
// It reports whether anything reached the display. A frame that redrew the
// same picture — a window sitting still with no animation running — sends
// nothing and answers false, and so costs almost nothing: no upload, no scan
// of the canvas, no drawing the backend had to be woken for.
func (win *Window) End() bool {
	win.moveFocus()
	// Drawings registered while the frame was being painted land here, above
	// everything, before the pixels are presented.
	for i := len(win.frontmost) - 1; i >= 0; i-- {
		win.frontmost[i]()
	}
	if win.native == nil {
		return false
	}
	dirty := win.dirtyRegion()
	// A frame that drew the picture already on the screen is not presented at
	// all: the platform is not asked, which is what keeps an idle window from
	// waking the compositor sixty times a second to be told there is nothing to
	// do. The canvas's own copy of the screen has already been brought up to
	// date by the measuring, so the next frame starts from the right place.
	if dirty.Width > 0 && dirty.Height > 0 {
		win.native.present(win, dirty)
		win.presented = true
	} else {
		win.presented = false
	}
	win.pace()
	return win.presented
}

// Presented reports whether the last [Window.End] put anything on the screen.
// A template can read it to see whether the frame it just drew was a change
// or a repeat: the input state is settled either way, and a repeat is free.
func (win *Window) Presented() bool { return win.presented }

// frameBudget is how long one frame may take. The default is the display's own
// refresh rate, so the loop is paced by the screen rather than by how fast the
// machine can go; a rate given to [Window.SetFPS] is a ceiling on top of that.
// A display that will not say is assumed to be the usual sixty.
func (win *Window) frameBudget() time.Duration {
	fps := win.targetFPS
	if fps <= 0 {
		if win.refreshHz == 0 && win.native != nil {
			win.refreshHz = win.native.displayRefresh()
		}
		if fps = win.refreshHz; fps <= 0 {
			fps = fallbackFPS
		}
	}
	return time.Duration(float64(time.Second) / float64(fps))
}

// pace holds the loop to the budget. The frames land on an absolute grid, so
// the time a frame took to draw does not add to the wait of the next one: a
// frame that overran its budget starts the next frame now and the loop drops
// frames rather than building a debt it can never pay back.
//
// The platform is told when the next frame is due, so a backend that waits on
// its own event queue can stay asleep until then rather than being woken to be
// told there is nothing to do.
func (win *Window) pace() {
	budget := win.frameBudget()
	if budget <= 0 {
		win.handOver(time.Time{})
		return
	}
	now := time.Now()
	if win.nextFrame.IsZero() {
		win.nextFrame = now
	}
	win.nextFrame = win.nextFrame.Add(budget)
	if !win.nextFrame.After(now) {
		win.nextFrame = now
		win.handOver(now)
		return
	}
	win.handOver(win.nextFrame)
	win.sleepUntil(win.nextFrame.Sub(now))
}

// handOver passes the frame deadline to a platform that can wait for it.
func (win *Window) handOver(due time.Time) {
	if pacer, ok := win.native.(framePacer); ok {
		pacer.setFrameDeadline(due)
	}
}

// sleepUntil waits out the rest of the frame. A platform with its own way to
// idle can be given one; the default is the obvious thing.
func (win *Window) sleepUntil(d time.Duration) {
	if win.sleep != nil {
		win.sleep(d)
		return
	}
	time.Sleep(d)
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

// SetFPS caps the frame rate. Zero or less — the default — follows the
// display, so a window is paced by the screen it is on; a positive number is a
// ceiling whatever the display says. A window that used to run flat out between
// Begin and End is now paced unless it asks otherwise.
func (win *Window) SetFPS(fps int) {
	win.targetFPS = max(fps, 0)
	// The old grid was set for the old rate; it means nothing now.
	win.nextFrame = time.Time{}
}

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

// SetOpacity fades the whole window — decorations and all — so that what sits
// behind the program shows through; 255 is fully opaque and 0 is invisible.
// It reports whether the request could be made: there is no window to fade on
// an offscreen canvas, so that reports false while still remembering the
// value. The window itself is never returned to — this is a whole-window
// opacity, not a fade of the pixels drawn into it, so what is drawn keeps its
// contrast and only its overall brightness changes.
func (win *Window) SetOpacity(alpha uint8) bool {
	win.opacity = alpha
	if win.native == nil {
		return false
	}
	return win.native.setOpacity(alpha)
}

// Opacity is what was last asked for with SetOpacity, 255 by default.
func (win *Window) Opacity() uint8 { return win.opacity }

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

// SetCanvas swaps the window's pixel surface for cv and returns the previous
// one, so a caller can paint into a scratch canvas and hand it back. The
// canvas must have the same size as the window — the width and height the
// window reports and the pixels it redraws stay pinned to that size — while
// the surface in between may hold a partially composed scene, like a layer
// waiting to be blitted over.
func (win *Window) SetCanvas(cv *canvas.Canvas) *canvas.Canvas {
	if cv == nil {
		return win.cv
	}
	prev := win.cv
	win.cv = cv
	return prev
}

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

// MouseMoved reports whether the pointer actually moved during this frame, as
// opposed to sitting still. A widget that lets the mouse move its selection
// (like a calendar's highlighted day) should only claim the selection on
// frames the pointer moved, so the keyboard keeps control otherwise.
func (win *Window) MouseMoved() bool { return win.mouseDX != 0 || win.mouseDY != 0 }

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

	// The copy of the last frame is the wrong size now, so it cannot be
	// compared against. The buffer behind it is kept: a window that is
	// resized once and then never again would otherwise throw away a
	// megabyte and take it again on the next frame.
	win.shadow = win.shadow[:0]
	win.fullRedraw = true
	win.safe = canvas.Area{}
	return true
}

// dirtyRegion is the rectangle to send: the pixels of this frame that really
// changed, which the canvas knows as each one is written because the canvas is
// comparing against the frame already on the screen. A frame that redraws the
// same picture changes nothing and sends nothing — there is no scan of the
// canvas left to find that out, because nothing has to look.
func (win *Window) dirtyRegion() canvas.Area {
	cv := win.cv
	whole := canvas.Area{X: 0, Y: 0, Width: cv.Width, Height: cv.Height}
	if win.fullRedraw || len(win.shadow) != cv.Width*cv.Height {
		win.fullRedraw = false
		win.adoptShadow(cv)
		return whole
	}
	a := cv.PresentChanges()
	return a
}

// adoptShadow takes the frame that is about to be sent whole as the one to
// compare against from now on, and gives the canvas the copy to check writes
// with. The pixels of the canvas are what the screen will hold, so the next
// frame starts from them.
func (win *Window) adoptShadow(cv *canvas.Canvas) {
	count := cv.Width * cv.Height
	if cap(win.shadow) < count {
		win.shadow = make([]canvas.Color, count)
	} else {
		win.shadow = win.shadow[:count]
	}
	copy(win.shadow, cv.Pixels)
	cv.Compare(win.shadow)
}
