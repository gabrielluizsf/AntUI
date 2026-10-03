package template

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// The suite draws into canvases it then looks at, and the loading screen
// draws over everything — so for every test that is not about the screen,
// the screen starts off here. The tests that ARE about it turn it back on
// around their own cases, either through loadingEnabledByDefault before a
// template is made or through the SetDisableLoading the program itself
// would use.
func init() { loadingEnabledByDefault = false }

// colorScreen paints the whole canvas one colour, standing in for the
// content of an app: what the loading screen is supposed to cover until
// the app is ready.
type colorScreen struct {
	c     canvas.Color
	drawn bool
}

func (s *colorScreen) Draw(c *Context) {
	c.Window().Canvas().Clear(s.c)
	s.drawn = true
}

// showsLoading says whether the canvas wears the loading screen: the
// background over everything and the bar across the middle. The bar is
// visible under the built-in template, whose track is the theme's border
// colour rather than the background it clears to.
func showsLoading(win *antui.Window, bg canvas.Color) bool {
	cv := win.Canvas()
	if cv.At(4, 4) != bg {
		return false
	}
	w := win.Width() * 3 / 5
	h := 6 * Scale(win)
	x0 := (win.Width() - w) / 2
	y := (win.Height()-h)/2 + h/2
	for x := x0; x < x0+w; x++ {
		if cv.At(x, y) != bg {
			return true
		}
	}
	return false
}

// showsBackground says whether the corner of the canvas is the plain
// background. The coord-free template draws a bar only when the stylesheet
// gives the progress role a colour of its own, so the tests of that
// template look for the clear rather than for the bar.
func showsBackground(win *antui.Window, bg canvas.Color) bool {
	return win.Canvas().At(4, 4) == bg
}

func TestLoadingIsOnByDefault(t *testing.T) {
	loadingEnabledByDefault = true
	defer func() { loadingEnabledByDefault = false }()

	l := newLoading()
	if !l.active() {
		t.Error("a new template starts with its loading screen on")
	}
	if got := l.getTimeout(); got != 5*time.Second {
		t.Errorf("the default timeout is %v, want the five seconds it is documented to be", got)
	}
}

func TestTheLoadingScreenComesWithTheContext(t *testing.T) {
	loadingEnabledByDefault = true
	defer func() { loadingEnabledByDefault = false }()

	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	bg := ctx.Template().Background(win)
	if !showsLoading(win, bg) {
		t.Fatal("the frame that comes with the context never landed on the canvas")
	}

	// Draw brings the screen back even with no screen pushed at all: a
	// program still loading has nothing to show, and the screen is the
	// showing.
	win.Canvas().Clear(canvas.Color(0))
	win.Begin()
	ctx.Draw()
	win.End()
	if !showsLoading(win, bg) {
		t.Error("Draw, with no screen pushed, did not draw the loading screen")
	}
}

func TestTheLoadingScreenStaysUntilReady(t *testing.T) {
	loadingEnabledByDefault = true
	defer func() { loadingEnabledByDefault = false }()

	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	bg := ctx.Template().Background(win)
	red := canvas.RGBA(220, 30, 30, 255)
	screen := &colorScreen{c: red}
	ctx.Go(screen)

	win.Begin()
	ctx.Draw()
	win.End()
	if !screen.drawn {
		t.Fatal("the app's screen was not drawn under the loading one")
	}
	if !showsLoading(win, bg) {
		t.Error("the loading screen did not cover the content")
	}

	ctx.Ready()
	win.Begin()
	ctx.Draw()
	win.End()
	if got := win.Canvas().At(4, 4); got != red {
		t.Errorf("after Ready the corner of the canvas is %v, want the screen's red %v", got, red)
	}
	if !win.Running() {
		t.Error("the window closed although the app was ready in time")
	}
}

func TestTheAppClosesWhenItsLoadingTimeoutPasses(t *testing.T) {
	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	ctx.SetDisableLoading(false)
	ctx.SetLoadingTimeout(60 * time.Millisecond)
	bg := ctx.Template().Background(win)

	start := time.Now()
	for win.Running() && time.Since(start) < 2*time.Second {
		win.Begin()
		ctx.Draw()
		win.End()
		time.Sleep(2 * time.Millisecond)
	}
	elapsed := time.Since(start)
	if win.Running() {
		t.Fatal("the app is still open past its loading timeout")
	}
	if elapsed < 50*time.Millisecond {
		t.Errorf("the app closed after %v, before the 60ms it was given", elapsed)
	}
	if elapsed > time.Second {
		t.Errorf("the app closed after %v: it waited for the default rather than for the timeout it was given", elapsed)
	}
	if ld := ctx.Template().(*uiTemplate).ld; !ld.timedOut.Load() {
		t.Error("the window closed without the loading timeout being what closed it")
	}
	if !showsLoading(win, bg) {
		t.Error("the last frame was not the loading screen with its bar at full")
	}
}

func TestAlongerLoadingTimeoutDelaysTheQuit(t *testing.T) {
	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	ctx.SetDisableLoading(false)
	ctx.SetLoadingTimeout(250 * time.Millisecond)

	start := time.Now()
	for time.Since(start) < 100*time.Millisecond {
		win.Begin()
		ctx.Draw()
		win.End()
		time.Sleep(2 * time.Millisecond)
	}
	if !win.Running() {
		t.Fatalf("the app closed after %v, well inside the 250ms it was given", time.Since(start))
	}
	for win.Running() && time.Since(start) < 2*time.Second {
		win.Begin()
		ctx.Draw()
		win.End()
		time.Sleep(2 * time.Millisecond)
	}
	if win.Running() {
		t.Fatal("the app never closed after its 250ms timeout")
	}
	if elapsed := time.Since(start); elapsed < 200*time.Millisecond {
		t.Errorf("the app closed after %v, before the 250ms it was given", elapsed)
	}
}

func TestANonPositiveLoadingTimeoutChangesNothing(t *testing.T) {
	l := newLoading()
	l.setTimeout(0)
	l.setTimeout(-time.Second)
	if got := l.getTimeout(); got != defaultLoadingTimeout {
		t.Errorf("after a zero and a negative duration the timeout is %v, want the %v it already had", got, defaultLoadingTimeout)
	}
	l.setTimeout(250 * time.Millisecond)
	if got := l.getTimeout(); got != 250*time.Millisecond {
		t.Errorf("the timeout is %v, want the 250ms it was set to", got)
	}
}

func TestTheProgressBarCountsUpToTheTimeout(t *testing.T) {
	l := newLoading()
	l.setTimeout(4 * time.Second)

	l.start = time.Now().Add(-2 * time.Second)
	if got := l.progress(); got < 0.49 || got > 0.51 {
		t.Errorf("halfway to the timeout the bar reads %v, want about 0.5", got)
	}

	l.start = time.Now().Add(-8 * time.Second)
	if got := l.progress(); got != 1 {
		t.Errorf("past the timeout the bar reads %v, want it at full", got)
	}

	if got := progressAt(0, 4*time.Second); got != 0 {
		t.Errorf("at the start of the wait the bar reads %v, want 0", got)
	}
	if got := progressAt(time.Second, 0); got != 0 {
		t.Errorf("with no timeout to count against the bar reads %v, want 0", got)
	}
}

func TestABareTemplateHasNoLoadingScreen(t *testing.T) {
	// A template built as a bare struct rather than through [New] — the
	// shape a test in this package builds one in — carries no screen, and
	// every call about one passes over it rather than falling over.
	ut := &uiTemplate{style: builtinStyle{}}
	if ut.loadingState() != nil {
		t.Fatal("a bare template answered for a loading screen it does not have")
	}
	ut.Ready()
	ut.SetDisableLoading(true)
	ut.SetLoadingTimeout(time.Second)
	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ut.loadingState().arm(win, ut)
}

func TestReadyCanComeFromAnotherGoroutine(t *testing.T) {
	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	ctx.SetDisableLoading(false)
	ctx.SetLoadingTimeout(400 * time.Millisecond)
	red := canvas.RGBA(220, 30, 30, 255)
	ctx.Go(&colorScreen{c: red})

	done := make(chan struct{})
	go func() {
		time.Sleep(30 * time.Millisecond)
		ctx.Ready()
		close(done)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		win.Begin()
		ctx.Draw()
		win.End()
		if win.Canvas().At(4, 4) == red {
			break
		}
		time.Sleep(2 * time.Millisecond)
	}
	<-done

	if win.Canvas().At(4, 4) != red {
		t.Fatal("the content never came up after Ready")
	}
	if !win.Running() {
		t.Error("the app closed although Ready arrived before the timeout")
	}
}

func TestADisabledLoadingScreenNeitherShowsNorCloses(t *testing.T) {
	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	bg := ctx.Template().Background(win)
	ctx.SetDisableLoading(false) // on for a moment...
	ctx.SetDisableLoading(true)  // ...and off, the way a program opts out
	ctx.SetLoadingTimeout(40 * time.Millisecond)

	deadline := time.Now().Add(250 * time.Millisecond)
	for time.Now().Before(deadline) {
		win.Begin()
		ctx.Draw()
		win.End()
		time.Sleep(2 * time.Millisecond)
	}
	if !win.Running() {
		t.Error("the app closed although the loading screen was disabled")
	}
	if showsBackground(win, bg) {
		t.Error("the disabled loading screen drew on the canvas anyway")
	}
	if ld := ctx.Template().(*uiTemplate).ld; ld.timedOut.Load() {
		t.Error("the loading timeout fired although the screen was disabled")
	}
}

func TestDisablingTheLoadingScreenWipesItFromTheCanvas(t *testing.T) {
	loadingEnabledByDefault = true
	defer func() { loadingEnabledByDefault = false }()

	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext(win, Builtin())
	bg := ctx.Template().Background(win)
	if !showsLoading(win, bg) {
		t.Fatal("no loading frame came with the context")
	}

	ctx.SetDisableLoading(true)
	if showsLoading(win, bg) {
		t.Error("the loading picture is still on the canvas after disabling")
	}
	if got := win.Canvas().At(4, 4); got != bg {
		t.Errorf("the wipe left %v on the canvas, want the plain background", got)
	}
}

func TestTheCSSLoopShowsItsLoadingScreenUntilReady(t *testing.T) {
	loadingEnabledByDefault = true
	defer func() { loadingEnabledByDefault = false }()

	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	bg := tpl.Background(win)
	if !showsBackground(win, bg) {
		t.Fatal("the loading frame did not come with the template")
	}

	// The loop of a coord-free app, with red standing in for its widgets.
	red := canvas.RGBA(220, 30, 30, 255)
	frame := func() {
		win.Begin()
		win.Clear(bg)
		tpl.Reset()
		win.FillRect(0, 0, win.Width(), win.Height(), red)
		win.End()
	}

	frame()
	if got := win.Canvas().At(4, 4); got != bg {
		t.Errorf("the loading screen did not cover the content: the corner is %v", got)
	}

	tpl.Ready()
	frame()
	if got := win.Canvas().At(4, 4); got != red {
		t.Errorf("after Ready the corner is %v, want the content's red %v", got, red)
	}
}

func TestTheCSSTemplateOptsOutOfLoading(t *testing.T) {
	loadingEnabledByDefault = true
	defer func() { loadingEnabledByDefault = false }()

	win, _, err := antui.Offscreen(240, 160)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	bg := tpl.Background(win)

	// Whatever is on the canvas when the flag goes down is replaced by the
	// plain background: the picture the screen left behind, wiped.
	marker := canvas.RGBA(220, 30, 30, 255)
	win.FillRect(0, 0, win.Width(), win.Height(), marker)
	tpl.SetDisableLoading(true)
	if got := win.Canvas().At(4, 4); got != bg {
		t.Errorf("disabling left %v on the canvas, want the plain background", got)
	}

	// And with the screen off, the timeout never fires.
	tpl.SetLoadingTimeout(40 * time.Millisecond)
	deadline := time.Now().Add(250 * time.Millisecond)
	for time.Now().Before(deadline) {
		win.Begin()
		win.Clear(bg)
		tpl.Reset()
		win.End()
		time.Sleep(2 * time.Millisecond)
	}
	if !win.Running() {
		t.Error("the app closed although the CSS template had the screen disabled")
	}
}
