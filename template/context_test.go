package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/template/event"
)

// screenSpy records whether Draw was called.
type screenSpy struct{ drawn bool }

func (s *screenSpy) Draw(c *Context) { s.drawn = true }

func TestContextPushAndDraw(t *testing.T) {
	win, _, _ := antui.Offscreen(300, 300)
	ctx := NewContext(win, Cyberpunk(win))
	sp := &screenSpy{}
	ctx.Go(sp)

	win.Begin()
	ctx.Draw()
	win.End()

	if !sp.drawn {
		t.Error("a pushed screen should be drawn on the next Draw")
	}
}

func TestContextBack(t *testing.T) {
	win, _, _ := antui.Offscreen(200, 200)
	ctx := NewContext(win, Builtin(win))
	first := &screenSpy{}
	ctx.Go(first)
	second := &screenSpy{}
	ctx.Go(second)

	if ctx.Current() != second {
		t.Fatal("current should be the last pushed screen")
	}

	ctx.Back()
	if ctx.Current() != first {
		t.Error("Back should pop to the screen underneath")
	}
	if ctx.Depth() != 1 {
		t.Error("Back should leave one screen on the stack")
	}

	// Back should not pop the last screen.
	ctx.Back()
	if ctx.Current() != first {
		t.Error("Back should not pop the first screen")
	}

	// Push three more, then Back once.
	ctx.Go(&screenSpy{})
	ctx.Go(&screenSpy{})
	ctx.Go(&screenSpy{})
	n := ctx.Depth()
	if n != 4 {
		t.Fatalf("Depth = %d, want 4", n)
	}
	ctx.Back()
	if ctx.Depth() != n-1 {
		t.Errorf("Depth after Back = %d, want %d", ctx.Depth(), n-1)
	}
}

func TestContextHome(t *testing.T) {
	win, _, _ := antui.Offscreen(200, 200)
	ctx := NewContext(win, Builtin(win))
	ctx.Go(&screenSpy{})
	ctx.Go(&screenSpy{})
	ctx.Go(&screenSpy{})
	ctx.Home()
	if ctx.Depth() != 1 {
		t.Errorf("Home should collapse to depth 1, got %d", ctx.Depth())
	}
}

func TestContextReplace(t *testing.T) {
	win, _, _ := antui.Offscreen(200, 200)
	ctx := NewContext(win, Builtin(win))
	ctx.Go(&screenSpy{})

	after := &screenSpy{}
	ctx.Replace(after)
	if ctx.Current() != after {
		t.Error("Replace should swap the current screen")
	}
	if ctx.Depth() != 1 {
		t.Error("Replace should not change the depth")
	}
}

type countingScreen struct{ n *int }

func (s *countingScreen) Draw(c *Context) { *s.n++ }

func TestContextDrawCallsScreen(t *testing.T) {
	win, _, _ := antui.Offscreen(128, 128)
	ctx := NewContext(win, Cyberpunk(win))
	var count int
	ctx.Go(&countingScreen{&count})
	win.Begin()
	ctx.Draw()
	win.End()
	if count != 1 {
		t.Errorf("Draw should call the screen once, drew %d", count)
	}
}

// -----------------------------------------------------------------------
// The navigation example: home → shop → cart, and back.
// -----------------------------------------------------------------------

type navHome struct{}
type navShop struct{}
type navCart struct{}

func (navHome) Draw(c *Context) {
	c.Label(0, 0, "Home")
	if c.Button(0, 40, 80, 30, "Shop").Is(event.Button, event.Click) {
		c.Go(navShop{})
	}
}

func (navShop) Draw(c *Context) {
	c.Label(0, 0, "Shop")
	c.BackButton()
	if c.Button(0, 40, 80, 30, "Cart").Is(event.Button, event.Click) {
		c.Go(navCart{})
	}
}

func (navCart) Draw(c *Context) {
	c.Label(0, 0, "Cart")
	c.BackButton()
	if c.Button(0, 40, 80, 30, "Home").Is(event.Button, event.Click) {
		c.Go(navHome{})
	}
}

func TestContextNavigationStack(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	ctx := NewContext(win, Cyberpunk(win))
	ctx.Go(navHome{})

	// Click "Shop".
	win.Begin()
	testClick(t, win, 40, 55)
	ctx.Draw()
	win.End()
	if ctx.Depth() != 2 {
		t.Fatalf("expected depth 2 after Go(Shop), got %d", ctx.Depth())
	}

	// Click "Cart".
	win.Begin()
	testClick(t, win, 40, 55)
	ctx.Draw()
	win.End()
	if ctx.Depth() != 3 {
		t.Fatalf("expected depth 3 after Go(Cart), got %d", ctx.Depth())
	}

	// Click BackButton (top-left area).
	win.Begin()
	testClick(t, win, 30, 25)
	ctx.Draw()
	win.End()
	if ctx.Depth() != 2 {
		t.Fatalf("expected depth 2 after Back from Cart, got %d", ctx.Depth())
	}
}

// testClick pushes a full press-and-release at a point. It must be called
// between Begin and End: Begin clears the per-frame input, then the pushes
// set the frames, and Button reads them.
func testClick(t *testing.T, win *antui.Window, x, y int) {
	t.Helper()
	win.Push(antui.Event{Type: antui.EventMouseMove, X: x, Y: y})
	win.Push(antui.Event{Type: antui.EventMouseDown, Button: antui.MouseLeft, X: x, Y: y})
	win.Push(antui.Event{Type: antui.EventMouseUp, Button: antui.MouseLeft, X: x, Y: y})
}
