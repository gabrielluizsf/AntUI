package antui

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// TestClaimedModalBlocksOtherWidgets guards a bug where a widget drawn
// earlier in flow order, underneath where a popup — a calendar, a dropdown
// menu — would land, still took a click that landed on it, because nothing
// about an open popup reached a widget that had already run its own hit
// test before the popup claimed the screen. A popup opened on an earlier
// frame and still open answers that for every widget in the current one,
// including those drawn before the popup's own call — and including one
// sitting directly underneath the popup's own rectangle, which is exactly
// where a widget the popup covers is.
func TestClaimedModalBlocksOtherWidgets(t *testing.T) {
	win, _, err := Offscreen(400, 400)
	if err != nil {
		t.Fatalf("offscreen: %v", err)
	}

	const modalOwner, otherWidget uint32 = 1, 2
	win.ClaimModal(modalOwner)

	// otherWidget sits inside the rectangle the popup itself occupies —
	// where the widget it is covering actually is — not off to the side of
	// it, which is the case most worth guarding against.
	if !win.Begin() {
		t.Fatal("window closed")
	}
	win.Push(Event{Type: EventMouseMove, X: 110, Y: 110})
	win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: 110, Y: 110})
	win.Push(Event{Type: EventMouseUp, Button: MouseLeft, X: 110, Y: 110})
	hovered := win.Hovered(canvas.Area{X: 100, Y: 100, Width: 100, Height: 100})
	if win.WidgetClick(otherWidget, hovered) {
		t.Error("a widget underneath the open popup was clicked through it")
	}
	win.End()

	// The modal's own widget is never blocked by its own claim.
	if !win.Begin() {
		t.Fatal("window closed")
	}
	win.Push(Event{Type: EventMouseMove, X: 110, Y: 110})
	win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: 110, Y: 110})
	win.Push(Event{Type: EventMouseUp, Button: MouseLeft, X: 110, Y: 110})
	hovered = win.Hovered(canvas.Area{X: 100, Y: 100, Width: 100, Height: 100})
	if !win.WidgetClick(modalOwner, hovered) {
		t.Error("the modal's own widget could not be clicked inside its own claimed area")
	}
	win.End()

	// Releasing the claim lets other widgets answer again.
	win.ReleaseModal(modalOwner)
	if !win.Begin() {
		t.Fatal("window closed")
	}
	win.Push(Event{Type: EventMouseMove, X: 110, Y: 110})
	win.Push(Event{Type: EventMouseDown, Button: MouseLeft, X: 110, Y: 110})
	win.Push(Event{Type: EventMouseUp, Button: MouseLeft, X: 110, Y: 110})
	hovered = win.Hovered(canvas.Area{X: 100, Y: 100, Width: 100, Height: 100})
	if !win.WidgetClick(otherWidget, hovered) {
		t.Error("a widget still could not be clicked after the modal released its claim")
	}
	win.End()
}
