package template

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// iconPainted is how many pixels of a rectangle a drawn icon covers, which is
// what tells "an icon is there" apart from "the box is empty".
func iconPainted(cv *canvas.Canvas, x, y, w, h int) int {
	n := 0
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if px < 0 || py < 0 || px >= cv.Width || py >= cv.Height {
				continue
			}
			if cv.At(px, py).A() != 0 {
				n++
			}
		}
	}
	return n
}

// mask is which pixels of a rectangle a mark covers, and the two ways a widget
// is painted a difference is measured between.
type mask struct {
	w, h  int
	marks []bool
}

func (m mask) at(x, y int) bool { return m.marks[y*m.w+x] }

func (m mask) count() (n int) {
	for _, on := range m.marks {
		if on {
			n++
		}
	}
	return n
}

// markOf paints a widget two ways and answers the pixels where the two differ.
//
// A mark has to be measured as a difference. A widget's own surface is opaque,
// so counting painted pixels counts the box; and counting the pixels of the
// mark's own colour counts the box as well whenever the two are near each other,
// which they are by default — a white tick on a white box. What the mark *is*
// is the difference between the widget with it and the same widget without.
func markOf(t *testing.T, w, h int, with, without func(win *antui.Window)) mask {
	t.Helper()
	paint := func(draw func(win *antui.Window)) *canvas.Canvas {
		win := blankWin(t, w, h)
		draw(win)
		return win.Canvas()
	}
	a, b := paint(with), paint(without)
	m := mask{w: w, h: h, marks: make([]bool, w*h)}
	for y := range h {
		for x := range w {
			m.marks[y*w+x] = a.At(x, y) != b.At(x, y)
		}
	}
	return m
}

// ends is the first and last row of a rectangle holding any of the mark, and
// the leftmost marked pixel on each. Which way a mark leans is read off it: a
// caret pointing down has its arms on the first row and its point on the last,
// so its leftmost pixel moves to the right between them.
func (m mask) ends() (firstRow, lastRow, firstLeft, lastLeft int, found bool) {
	for y := range m.h {
		for x := range m.w {
			if !m.at(x, y) {
				continue
			}
			if !found {
				firstRow, firstLeft, found = y, x, true
			}
			lastRow, lastLeft = y, x
			break
		}
	}
	return firstRow, lastRow, firstLeft, lastLeft, found
}

// extent is the box the mark covers, which is what says whether a mark is in
// the middle of a control or off in one corner of it.
func (m mask) extent() (minX, minY, maxX, maxY int, found bool) {
	// The extremes have to be tracked one at a time rather than taken from the
	// first and last pixels to come up: in a round mark the first pixel is on a
	// short row, and the widest row is in the middle.
	minX, minY, maxX, maxY = m.w, m.h, -1, -1
	for i, on := range m.marks {
		if !on {
			continue
		}
		found = true
		x, y := i%m.w, i/m.w
		minX, maxX = min(minX, x), max(maxX, x)
		minY, maxY = min(minY, y), max(maxY, y)
	}
	return minX, minY, maxX, maxY, found
}

// restoreIcons puts the built-in icons back after a test that changed one, so
// that a case which replaces an icon does not decide what the next one draws.
func restoreIcons(t *testing.T) {
	t.Helper()
	t.Cleanup(resetIcons)
}

func TestBuiltinIconsAreRegisteredAndDraw(t *testing.T) {
	// Every icon the widgets draw has to be registered and to paint something:
	// an icon that is registered but empty is a widget with a hole in it.
	win := blankWin(t, 120, 30)
	cv := win.Canvas()
	for i, name := range []string{IconCheck, IconDot, IconCaret, IconCaretUp} {
		if _, ok := IconSource(name); !ok {
			t.Errorf("the built-in icon %q is not registered", name)
			continue
		}
		if !DrawIcon(cv, name, i*30, 0, 24, 24, canvas.RGB(255, 255, 255)) {
			t.Errorf("the built-in icon %q did not draw", name)
			continue
		}
		if got := iconPainted(cv, i*30, 0, 24, 24); got < 20 {
			t.Errorf("the built-in icon %q painted %d pixels, want a visible mark", name, got)
		}
	}
}

func TestIconSourceIsTheDrawingToChange(t *testing.T) {
	// The text an icon is registered from is the whole point of registering it
	// as SVG: a developer reads it, changes it, and hands it back. Reading one
	// and getting the file back unchanged would make that impossible.
	restoreIcons(t)
	src, ok := IconSource(IconCheck)
	if !ok {
		t.Fatal("the tick is not registered")
	}
	if !strings.Contains(src, "currentColor") {
		t.Errorf("the built-in tick does not follow the colour, so it could not serve several states:\n%s", src)
	}
	if !strings.Contains(src, "<svg") {
		t.Errorf("the source is not a drawing:\n%s", src)
	}
}

func TestSetIconChangesWhatEveryWidgetDraws(t *testing.T) {
	// A changed drawing replaces the mark everywhere it is used: the widget asks
	// for the icon by name, so one change reaches the checkbox of every list
	// without any of them being told about it.
	restoreIcons(t)
	win := blankWin(t, 40, 40)
	cv := win.Canvas()
	before := iconPainted(cv, 0, 0, 0, 0)

	// A tick that is a filled square instead of a stroke: a shape a different
	// icon set would have, and one no recolouring could produce.
	if err := SetIcon(IconCheck, `<svg viewBox="0 0 24 24"><rect width="24" height="24" fill="currentColor"/></svg>`); err != nil {
		t.Fatalf("SetIcon: %v", err)
	}
	_ = before
	DrawIcon(cv, IconCheck, 0, 0, 24, 24, canvas.RGB(255, 255, 255))
	if got := iconPainted(cv, 0, 0, 24, 24); got < 24*24-200 {
		t.Errorf("the replaced tick painted %d pixels, want a filled square", got)
	}
	// The other icons are untouched by the change.
	DrawIcon(cv, IconCaret, 0, 0, 0, 0, canvas.RGB(255, 255, 255))
}

func TestSetIconReportsADrawingItCannotRead(t *testing.T) {
	// A drawing that cannot be read is reported to whoever registered it, and
	// leaves the icon that was there: a broken icon should not empty a widget
	// that was drawing a moment ago.
	restoreIcons(t)
	if err := SetIcon(IconCheck, `<svg><rect></svg>`); err == nil {
		t.Fatal("a drawing that does not close was accepted")
	}
	if _, ok := IconSource(IconCheck); !ok {
		t.Error("the tick is gone after a failed change")
	}
	if err := SetIcon("  ", `<svg/>`); err == nil {
		t.Error("an icon with no name was accepted")
	}
}

func TestSetIconRegistersANewName(t *testing.T) {
	// A name nobody has used is a new icon rather than a replacement, which is
	// how a button comes to wear one of the application's own.
	restoreIcons(t)
	if err := SetIcon("save", `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8" fill="currentColor"/></svg>`); err != nil {
		t.Fatalf("SetIcon: %v", err)
	}
	if _, ok := IconSource("save"); !ok {
		t.Error("the new icon is not registered under its own name")
	}
	found := false
	for _, name := range IconNames() {
		if name == "save" {
			found = true
		}
	}
	if !found {
		t.Errorf("IconNames does not list the new icon: %v", IconNames())
	}
	win := blankWin(t, 30, 30)
	if !DrawIcon(win.Canvas(), "save", 0, 0, 24, 24, canvas.RGB(255, 255, 255)) {
		t.Error("the new icon did not draw")
	}
}

func TestIconNamesAreSorted(t *testing.T) {
	// A tool that lets a developer pick an icon shows them in an order; the order
	// a map is walked in is not one.
	names := IconNames()
	if len(names) < 4 {
		t.Fatalf("only %d icons are registered, want the built-in ones", len(names))
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Fatalf("IconNames is not sorted: %v", names)
		}
	}
}

func TestIconCanvasIsRememberedAndRepaintedWhenItShouldBe(t *testing.T) {
	// The picture of an icon at a size in a colour is drawn once and answered
	// again, and drawn afresh when either of those changes — the second half is
	// what a list of buttons where one is hovered depends on.
	first := IconCanvas(IconCheck, 24, 24, canvas.RGB(255, 0, 0))
	if first == nil {
		t.Fatal("the tick did not paint")
	}
	if again := IconCanvas(IconCheck, 24, 24, canvas.RGB(255, 0, 0)); again != first {
		t.Error("asking twice for the same icon drew it twice")
	}
	if other := IconCanvas(IconCheck, 24, 24, canvas.RGB(0, 0, 255)); other == first {
		t.Error("a second colour answered the canvas of the first")
	}
	if bigger := IconCanvas(IconCheck, 48, 48, canvas.RGB(255, 0, 0)); bigger == nil || bigger.Width != 48 {
		t.Error("a second size answered the canvas of the first")
	}
}

func TestDrawIconOfNothingIsNotAnError(t *testing.T) {
	// A widget may ask for an icon that was never registered — one a
	// stylesheet named and nothing supplied — and the widget still has to draw
	// the rest of itself.
	win := blankWin(t, 20, 20)
	cv := win.Canvas()
	if DrawIcon(cv, "there-is-no-such-icon", 0, 0, 16, 16, canvas.White) {
		t.Error("an icon that is not registered drew something")
	}
	if DrawIcon(nil, IconCheck, 0, 0, 16, 16, canvas.White) {
		t.Error("drawing on nothing reported a drawing")
	}
	if DrawIcon(cv, IconCheck, 0, 0, 0, 16, canvas.White) {
		t.Error("a size of nothing reported a drawing")
	}
	if got := iconPainted(cv, 0, 0, 20, 20); got != 20*20 {
		t.Errorf("something was painted where no icon was asked for: %d pixels", got)
	}
}

func TestWidgetsDrawTheirMarksAsIcons(t *testing.T) {
	// The widgets that used to draw a mark with straight lines now draw the icon
	// of that mark. The mark is measured as a difference between the same control
	// painted with its icon and with that icon emptied, which is the mark itself
	// and nothing else — a control that is on and one that is off differ in their
	// fill and their border as well, so comparing the two states would measure
	// the box rather than the tick.
	restoreIcons(t)
	const box = 18 // the box is 18 units at the scale of this window, which is one
	const blankDrawing = `<svg viewBox="0 0 24 24"></svg>`

	markOfIcon := func(draw func(win *antui.Window), name string) mask {
		return markOf(t, box, box,
			func(win *antui.Window) {
				resetIcons()
				draw(win)
			},
			func(win *antui.Window) {
				if err := SetIcon(name, blankDrawing); err != nil {
					t.Fatalf("SetIcon: %v", err)
				}
				draw(win)
			})
	}

	tick := markOfIcon(func(win *antui.Window) {
		builtinStyle{}.Checkbox(win, State{}, 0, 0, "", true)
	}, IconCheck)
	if tick.count() < 8 {
		t.Errorf("a checked checkbox has %d pixels of tick in it, want a visible mark", tick.count())
	}
	if _, lastRow, firstLeft, _, ok := tick.ends(); !ok {
		t.Error("the tick is not in the box")
	} else {
		// Inside the box, and clear of its border.
		if lastRow >= box {
			t.Errorf("the tick reaches row %d of a %d box, want it inside", lastRow, box)
		}
		if firstLeft == 0 {
			t.Error("the tick starts at the very edge of the box, want room for the border")
		}
	}
	// An unchecked box shows no tick: the mark belongs to being on.
	if off := markOfIcon(func(win *antui.Window) {
		builtinStyle{}.Checkbox(win, State{}, 0, 0, "", false)
	}, IconCheck); off.count() != 0 {
		t.Errorf("an unchecked checkbox has %d pixels of tick in it, want none", off.count())
	}

	dot := markOfIcon(func(win *antui.Window) {
		builtinStyle{}.Radio(win, State{}, 0, 0, "", true)
	}, IconDot)
	if dot.count() < 8 {
		t.Errorf("a selected radio has %d pixels of dot in it, want a visible mark", dot.count())
	}
	// The dot is in the middle of the box, which is what a radio's is, and it is
	// round rather than spread across the whole thing.
	minX, minY, maxX, maxY, ok := dot.extent()
	if !ok {
		t.Fatal("the selected radio has no dot in it")
	}
	if mid := (box - 1) / 2; (minX+maxX)/2 < mid-2 || (minX+maxX)/2 > mid+2 {
		t.Errorf("the dot sits at columns %d..%d of a %d box, want it in the middle", minX, maxX, box)
	}
	if w, h := maxX-minX+1, maxY-minY+1; w > box/2 || h > box/2 {
		t.Errorf("the dot is %dx%d in a %d box, want a dot rather than a mark across it", w, h, box)
	}
	t.Logf("extent: minX=%d minY=%d maxX=%d maxY=%d found=%v box=%d  mask.w=%d mask.h=%d len=%d", minX, minY, maxX, maxY, ok, box, dot.w, dot.h, len(dot.marks))
	if w, h := maxX-minX+1, maxY-minY+1; w < 3 || h < 3 || w != h {
		t.Errorf("the dot is %dx%d, want something round of a few pixels across", w, h)
	}
	if off := markOfIcon(func(win *antui.Window) {
		builtinStyle{}.Radio(win, State{}, 0, 0, "", false)
	}, IconDot); off.count() != 0 {
		t.Errorf("an unselected radio has %d pixels of dot in it, want none", off.count())
	}
}

func TestSelectAndDatePickerDrawTheirCarets(t *testing.T) {
	// A dropdown and a date picker both end in an arrow saying there is more
	// below, and the date picker's turns over while its calendar is open. The
	// arrow is found by putting an icon that paints nothing in its place and
	// keeping the difference — which is also the proof that the icon is what
	// draws the arrow, and not some line left behind.
	restoreIcons(t)
	const w, h = 140, 26
	const blankDrawing = `<svg viewBox="0 0 24 24"></svg>`
	blank := func(name string) {
		t.Helper()
		if err := SetIcon(name, blankDrawing); err != nil {
			t.Fatalf("SetIcon: %v", err)
		}
	}

	// Two paints of the same widget with the same icons have to agree, or the
	// difference measured below would be measuring nothing at all.
	if diff := markOf(t, w, h,
		func(win *antui.Window) { builtinStyle{}.Select(win, State{}, 0, 0, w, h, "v", false) },
		func(win *antui.Window) { builtinStyle{}.Select(win, State{}, 0, 0, w, h, "v", false) }); diff.count() != 0 {
		t.Fatalf("a widget painted twice differed in %d pixels, so a mark cannot be told from its widget", diff.count())
	}

	arrow := markOf(t, w, h,
		func(win *antui.Window) {
			resetIcons()
			builtinStyle{}.Select(win, State{}, 0, 0, w, h, "v", false)
		},
		func(win *antui.Window) {
			blank(IconCaret)
			builtinStyle{}.Select(win, State{}, 0, 0, w, h, "v", false)
		})
	if arrow.count() == 0 {
		t.Fatal("the dropdown has no arrow in it")
	}
	// The arrow is at the right-hand end of the box, which is where a dropdown's
	// is, and clear of the value at the left.
	firstRow, lastRow, firstLeft, _, _ := arrow.ends()
	if firstLeft < w/2 {
		t.Errorf("the arrow starts at column %d of a %d box, want it at the right-hand end", firstLeft, w)
	}
	if lastRow-firstRow < 2 {
		t.Errorf("the arrow is %d rows tall, want a mark rather than a line", lastRow-firstRow+1)
	}

	// The date picker's arrow is the same one, and it is the other way up while
	// the calendar is open. A caret pointing down has its two arms at the top and
	// its point at the bottom, so its leftmost pixel moves to the right going
	// down; one pointing up is the other way round.
	caret := func(name string, open bool) mask {
		return markOf(t, w, h,
			func(win *antui.Window) {
				resetIcons()
				builtinStyle{}.DatePickerBox(win, State{}, 0, 0, w, h, "v", open)
			},
			func(win *antui.Window) {
				blank(name)
				builtinStyle{}.DatePickerBox(win, State{}, 0, 0, w, h, "v", open)
			})
	}
	_, _, downFirst, downLast, ok := caret(IconCaret, false).ends()
	if !ok {
		t.Fatal("the closed date picker has no arrow in it")
	}
	_, _, upFirst, upLast, ok := caret(IconCaretUp, true).ends()
	if !ok {
		t.Fatal("the open date picker has no arrow in it")
	}
	if downLast <= downFirst {
		t.Error("the closed date picker's arrow does not point down")
	}
	if upLast >= upFirst {
		t.Error("the open date picker's arrow is the same way up as the closed one")
	}
}

func TestStylesheetURLPaintsARegisteredIconInTheTextsColour(t *testing.T) {
	// This is the point of the whole thing: a stylesheet can put one of the
	// icons on a widget, and the icon takes the colour of that widget's own
	// text. So a button's icon follows the button's colour, and a box whose
	// colour is changed by a rule gets its icon changed with it, without the icon
	// being written again.
	restoreIcons(t)
	win := blankWin(t, 60, 60)
	cs := newCSSStyle(win)

	red := bgRule(t, "background-image: url(check); background-repeat: no-repeat; color: #ff0000;")
	blue := bgRule(t, "background-image: url(check); background-repeat: no-repeat; color: #0000ff;")

	paint := func(st css.Style) canvas.Canvas {
		cs.paintBackground(win, st, 0, 0, 60, 60, 0, 0)
		return *win.Canvas()
	}
	first := paint(red)
	redPixel := iconPainted(win.Canvas(), 0, 0, 60, 60)
	if redPixel == 0 {
		t.Fatal("a url() naming a registered icon painted nothing")
	}
	// Somewhere in what was painted is the red of the icon, not the background.
	found := false
	for y := 0; y < 60 && !found; y++ {
		for x := 0; x < 60; x++ {
			c := first.At(x, y)
			if c.R() > 200 && c.G() < 80 && c.B() < 80 {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("the icon did not paint in the colour the rule declared for the text")
	}

	// The same icon on a box whose text is another colour comes out that colour,
	// which is the whole reason the icon is painted per box and not once.
	win.Canvas().FillRect(0, 0, 60, 60, canvas.RGBA(30, 40, 50, 255))
	paint(blue)
	found = false
	for y := 0; y < 60 && !found; y++ {
		for x := 0; x < 60; x++ {
			c := win.Canvas().At(x, y)
			if c.B() > 200 && c.R() < 80 && c.G() < 80 {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("the icon kept the colour of the previous box instead of this one's text")
	}
}

func TestStylesheetCanvasWinsOverAnIconOfTheSameName(t *testing.T) {
	// A canvas registered on the template for a name that is also an icon is what
	// the stylesheet gets: it was named for this template on purpose, where an
	// icon is only what an unclaimed name falls back to.
	win := blankWin(t, 40, 40)
	cs := newCSSStyle(win)
	st := bgRule(t, "background-image: url(dot); background-repeat: no-repeat;")
	img, err := canvas.NewCanvas(12, 12)
	if err != nil {
		t.Fatal(err)
	}
	img.FillRect(0, 0, 12, 12, canvas.Red)
	cs.images["dot"] = img

	cs.paintBackground(win, st, 10, 10, 20, 20, 0, 0)
	if got := win.Canvas().At(15, 15); got != canvas.Red {
		t.Errorf("the registered canvas did not paint, got %v", got)
	}
}

func TestStylesheetIconFollowsTheWidgetScale(t *testing.T) {
	// An icon on a box is painted at the size it was drawn to be times the scale
	// of the window, so the same stylesheet gives the same proportions on a phone
	// and on a 4K screen, which is the rule every other measurement follows.
	// The two windows are one scale apart — 720 across is twice ScaleBase, 360 is
	// once — and the icon is measured on the window rather than asked for, since
	// what the widget drew is the thing being checked.
	paintedAt := func(width, height int) int {
		win := blankWin(t, width, height)
		cs := newCSSStyle(win)
		st := bgRule(t, "background-image: url(caret); background-repeat: no-repeat; color: #ffffff;")
		cs.paintBackground(win, st, 0, 0, width, height, 0, 0)
		return iconPainted(win.Canvas(), 0, 0, width, height)
	}

	small := paintedAt(ScaleBase, ScaleBase)
	large := paintedAt(ScaleBase*2, ScaleBase*2)
	if small == 0 {
		t.Fatalf("the caret was not painted on a %d wide window", ScaleBase)
	}
	// The icon is a caret, so most of the window is not icon: what is being
	// compared is the icon's area, which scales with the square of the scale.
	// Doubling the window quadruples it, less a little for the soft edge that
	// rounds off differently at each size.
	if large < small*3 {
		t.Errorf("the caret covered %d pixels at twice the scale and %d at once, want it to grow with the window", large, small)
	}
	if large > small*5 {
		t.Errorf("the caret covered %d pixels at twice the scale and %d at once, want it to grow with the window, not faster", large, small)
	}

	// The icon itself is asked for at the scale of the window, and comes back
	// that size.
	win := blankWin(t, ScaleBase*2, ScaleBase*2)
	if got := Scale(win); got != 2 {
		t.Fatalf("a %d wide window has scale %d, want 2", ScaleBase*2, got)
	}
	icon := IconCanvas("caret", 24*Scale(win), 24*Scale(win), canvas.White)
	if icon == nil {
		t.Fatal("the caret did not paint at the scale of the window")
	}
	if icon.Width != 48 {
		t.Errorf("the caret is %d wide on a doubled window, want 48", icon.Width)
	}
}

func TestIconsAreSafeFromSeveralGoroutines(t *testing.T) {
	// A program owns its icons while its windows paint on their own goroutines,
	// so the registry and the pictures are read from more than one at a time.
	// Each goroutine gets a canvas of its own: a canvas belongs to one window
	// and takes writes from one at a time, which is true of every drawing in the
	// library and not something the registry can promise to change.
	done := make(chan bool, 8)
	for i := range 8 {
		go func(i int) {
			cv := blankWin(t, 60, 60).Canvas()
			name := []string{IconCheck, IconCaret, IconDot, IconCaretUp}[i%4]
			for range 20 {
				DrawIcon(cv, name, i*7, 0, 16, 16, canvas.RGB(255, uint8(i*30), 0))
				IconNames()
				_, _ = IconSource(name)
			}
			done <- true
		}(i)
	}
	for range 8 {
		<-done
	}
}
