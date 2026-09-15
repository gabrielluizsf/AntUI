package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/event"
)

// A Screen is one view in a window: it draws itself on a Context and decides
// on its own when to leave. A screen knows nothing about the screens around
// it — the Context holds them, as a stack with the showing screen on top —
// so screens serve as the pointers a game navigates along:
//
//	home → shop → cart, and back the same way.
type Screen interface {
	Draw(c *Context)
}

// Context is one window running one template through a stack of screens. Screens
// push one another with Go, come back with Back, and the Context keeps the
// record — the stack IS the map home.
//
//	win, _ := antui.OpenWith(antui.Options{Title: "Shop", Width: 480, Height: 320})
//
//	app := template.NewContext(win, template.Cyberpunk(win))
//	app.Run(Home{})
//
// while Home's Draw decides when to leave:
//
//	func (Home) Draw(c *Context) {
//	    c.Label(0, 0, "Home")
//	    if c.Button(c.Center(160), 120, 160, 40, "To Shop").Is(event.Button, event.Click) {
//	        c.Go(Shop{})
//	    }
//	}
type Context struct {
	win   *antui.Window
	tpl   Template
	stack []Screen
}

// NewContext starts a session on a window with a template and nothing in
// view. Run, or Draw once a frame, shows the first screen.
func NewContext(win *antui.Window, tpl Template) *Context {
	return &Context{win: win, tpl: tpl}
}

// Go puts a screen on top of the stack. It becomes the current screen on the
// next [Context.Draw].
func (c *Context) Go(s Screen) {
	if s != nil {
		c.stack = append(c.stack, s)
	}
}

// Back pops the top screen off and shows the one underneath. The first screen
// ever shown stays: Back does nothing there.
func (c *Context) Back() {
	if len(c.stack) > 1 {
		c.stack = c.stack[:len(c.stack)-1]
	}
}

// Home collapses the stack back to the first screen that was shown, wherever
// it is. It is the "start over" of a navigation stack.
func (c *Context) Home() {
	if len(c.stack) > 1 {
		c.stack = c.stack[:1]
	}
}

// Replace swaps the top screen for another without going back first — the
// move a "logged in" screen makes after an entry form, so the form is not
// sitting underneath it waiting to be returned to.
func (c *Context) Replace(s Screen) {
	if s == nil {
		return
	}
	if len(c.stack) == 0 {
		c.stack = append(c.stack, s)
		return
	}
	c.stack[len(c.stack)-1] = s
}

// Current is the screen showing right now, or nil when none has been shown yet.
func (c *Context) Current() Screen {
	if len(c.stack) == 0 {
		return nil
	}
	return c.stack[len(c.stack)-1]
}

// Depth is how many screens are stacked, the current one included.
func (c *Context) Depth() int { return len(c.stack) }

// Draw paints the current screen once, cleared to the template's background
// first. Call it every frame between Begin and End; [Context.Run] does that
// for the whole life of the window.
func (c *Context) Draw() {
	if c.win == nil || len(c.stack) == 0 {
		return
	}
	c.win.Canvas().Clear(c.tpl.Background(c.win))
	c.stack[len(c.stack)-1].Draw(c)
}

// Run is the whole application loop: it shows the first screen and keeps
// drawing it (and whatever it pushes) until the window closes.
func (c *Context) Run(first Screen) {
	c.Go(first)
	for c.win.Begin() {
		c.Draw()
		c.win.End()
		if !c.win.Running() {
			break
		}
	}
}

// Template is the template the screens are drawn with.
func (c *Context) Template() Template { return c.tpl }

// Window is the window the screens draw on.
func (c *Context) Window() *antui.Window { return c.win }

// Width and Height are the window's drawing size, for centering and for
// laying a screen out against its edges.
func (c *Context) Width() int  { return c.win.Width() }
func (c *Context) Height() int { return c.win.Height() }

// Center is the left edge of a widget of the given width, centered
// horizontally — where a button goes to sit in the middle of the window.
func (c *Context) Center(w int) int { return (c.win.Width() - w) / 2 }

// Label draws a line of text in the template's own voice.
func (c *Context) Label(x, y int, text string) { c.tpl.Label(c.win, x, y, text) }

// Button, Checkbox, Radio, Slider and Input draw the template's components.
// Each reports an [event.Event]; screens use them to decide what happened.
func (c *Context) Button(x, y, w, h int, label string) event.Event {
	return c.tpl.Button(c.win, x, y, w, h, label)
}
func (c *Context) Checkbox(x, y int, label string, value *bool) event.Event {
	return c.tpl.Checkbox(c.win, x, y, label, value)
}
func (c *Context) Radio(x, y int, label string, value *int, option int) event.Event {
	return c.tpl.Radio(c.win, x, y, label, value, option)
}
func (c *Context) Slider(x, y, w, h int, value *float32, minValue, maxValue float32) event.Event {
	return c.tpl.Slider(c.win, x, y, w, h, value, minValue, maxValue)
}
func (c *Context) Input(x, y, w, h int, text *string) event.Event {
	return c.tpl.Input(c.win, x, y, w, h, text)
}

// BackButton draws a small "Back" button in the top-left corner that pops the
// screen when clicked, and reports whether it was clicked. It is the usual
// way out of a screen that was pushed on top of another. It scales with the
// window like everything else.
func (c *Context) BackButton() event.Event {
	u := c.Scale()
	w := (24 + canvas.TextWidth("Back")) * u
	e := c.tpl.Button(c.win, 12*u, 12*u, w, 26*u, "Back")
	if e.Is(event.Button, event.Click) {
		c.Back()
	}
	return e
}
