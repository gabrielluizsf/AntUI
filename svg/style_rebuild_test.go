package svg

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func TestRenderRebuildsDashWhenStyleChanges(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<line x1="0" y1="50" x2="100" y2="50" stroke="red" stroke-width="4"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	line := img.Root.find("line")
	first := img.Render(100, 100)
	if first == nil || pixelAt(first, 15, 50).A() == 0 {
		t.Fatal("the initial solid line was not painted")
	}
	dash, ok := canvas.NewDash([]float64{10, 10}, 0)
	if !ok {
		t.Fatal("could not make the test dash")
	}
	line.Style.Dash = &dash
	second := img.Render(100, 100)
	if second == nil || second == first {
		t.Fatal("changing Style.Dash did not invalidate the rendered canvas")
	}
	if got := pixelAt(second, 15, 50).A(); got != 0 {
		t.Errorf("the changed dash left the gap painted with alpha %d", got)
	}
}

func TestRenderRebuildsPathLengthWhenStyleChanges(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<line x1="0" y1="50" x2="100" y2="50" stroke="red" stroke-width="4"
			stroke-dasharray="4 2" pathLength="10"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	line := img.Root.find("line")
	img.Render(100, 100)
	if got := line.Style.Dash.On[0]; got != 40 {
		t.Fatalf("the parsed pathLength made the dash %g, want 40", got)
	}
	line.Style.pathLength = 20
	img.Render(100, 100)
	if got := line.Style.Dash.On[0]; got != 20 {
		t.Errorf("the changed pathLength made the dash %g, want 20", got)
	}
}

func TestRenderRebuildsMarkersWhenStyleChanges(t *testing.T) {
	img, err := Parse(`<svg viewBox="0 0 20 20">
		<defs>
			<marker id="red" markerUnits="userSpaceOnUse"><rect width="4" height="4" fill="red"/></marker>
			<marker id="blue" markerUnits="userSpaceOnUse"><rect width="4" height="4" fill="blue"/></marker>
		</defs>
		<line x1="10" y1="10" x2="18" y2="10" stroke="none" marker-start="url(#red)"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	line := img.Root.find("line")
	first := img.Render(20, 20)
	if got := pixelAt(first, 10, 10); got != canvas.RGB(255, 0, 0) {
		t.Fatalf("the initial marker is %v, want red", got)
	}
	line.Style.markerStart = img.markers["blue"]
	second := img.Render(20, 20)
	if second == first {
		t.Fatal("changing a marker did not invalidate the rendered canvas")
	}
	if got := pixelAt(second, 10, 10); got != canvas.RGB(0, 0, 255) {
		t.Errorf("the changed marker is %v, want blue", got)
	}
	if len(line.markers) != 1 || line.markers[0].def != img.markers["blue"] {
		t.Error("the marker placement was not rebuilt")
	}
}
