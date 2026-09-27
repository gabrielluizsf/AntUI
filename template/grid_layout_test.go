package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
	"github.com/gabrielluizsf/antui/template/event"
)

func gridRun(t *testing.T, sheet string, draw func(*CSS, *[]box)) []box {
	t.Helper()
	return gridRunIn(t, 400, 200, sheet, draw)
}

// gridRunIn draws the grid on a window of the given size, for the tracks whose
// count depends on how much room the container has.
func gridRunIn(t *testing.T, width, height int, sheet string, draw func(*CSS, *[]box)) []box {
	t.Helper()
	win, _, err := antui.Offscreen(width, height)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, sheet)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Grid(func(f *CSS) {
		draw(f, &got)
	})
	win.End()
	return got
}

func gridLayout(dst *[]box) func(*CSS, string) {
	return func(f *CSS, label string) {
		x, y, w, h, _, ok := f.layout(css.RoleButton, label)
		if ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	}
}

func assertGridBox(t *testing.T, got box, x, y, w, h int) {
	t.Helper()
	if got.x != x || got.y != y || got.w != w || got.h != h {
		t.Errorf("box = %d,%d %dx%d, want %d,%d %dx%d", got.x, got.y, got.w, got.h, x, y, w, h)
	}
}

// gridAxisOf is a resolved axis of n auto tracks, the shape the placement tests
// read lines off.
func gridAxisOf(n int) css.GridAxis {
	tracks := make([]css.GridTrack, n)
	for i := range tracks {
		tracks[i] = css.GridTrack{Kind: css.GridTrackSingle, Size: css.GridTrackSize{Kind: css.GridTrackAuto}}
	}
	return css.GridAxis{Tracks: tracks, Names: make([][]string, n+1)}
}

func TestGridAxisPlacementNormalizesLines(t *testing.T) {
	placement := css.GridPlacement{
		Start: css.GridLine{Kind: css.GridLineIndex, Index: 3},
		End:   css.GridLine{Kind: css.GridLineIndex, Index: 1},
	}
	start, end, ok := gridAxisPlacement(gridAxisOf(3), placement)
	if !ok || start != 0 || end != 2 {
		t.Errorf("reversed placement = %d,%d,%v", start, end, ok)
	}
	span := css.GridPlacement{
		Start: css.GridLine{Kind: css.GridLineSpan, Index: 1},
		End:   css.GridLine{Kind: css.GridLineSpan, Index: 2},
	}
	if got := gridPlacementSpan(span); got != 1 {
		t.Errorf("placement span = %d, want 1", got)
	}
	negative := css.GridPlacement{
		Start: css.GridLine{Kind: css.GridLineIndex, Index: -4},
		End:   css.GridLine{Kind: css.GridLineIndex, Index: -3},
	}
	if start, end, ok := gridAxisPlacement(gridAxisOf(2), negative); !ok || start != -1 || end != 0 {
		t.Errorf("leading negative placement = %d,%d,%v, want -1,0,true", start, end, ok)
	}
}

func TestGridAxisPlacementReadsLineNames(t *testing.T) {
	sheet, err := css.Parse(`grid { grid-template-columns: [full-start] 100px [content-start] 1fr [content-end] 100px [full-end]; }`)
	if err != nil {
		t.Fatal(err)
	}
	st := sheet.Style(css.RoleGrid, nil, css.StateNone, 400)
	axis := st.GridTemplateColumns.Resolve(400, 0, 1)
	for _, tc := range []struct {
		placement css.GridPlacement
		start     int
		end       int
		ok        bool
	}{
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "full-start"}}, 0, 1, true},
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "content-start"}}, 1, 2, true},
		{css.GridPlacement{
			Start: css.GridLine{Kind: css.GridLineName, Name: "content-start"},
			End:   css.GridLine{Kind: css.GridLineName, Name: "content-end"},
		}, 1, 2, true},
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "full-end"}}, 3, 4, true},
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "full-end", Backward: true}}, 3, 4, true},
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "content-start", Span: 2}}, 1, 3, true},
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "nowhere"}}, 0, 0, false},
		{css.GridPlacement{Start: css.GridLine{Kind: css.GridLineName, Name: "Full-Start"}}, 0, 0, false},
	} {
		start, end, got := gridAxisPlacement(axis, tc.placement)
		if start != tc.start || end != tc.end || got != tc.ok {
			t.Errorf("%+v = %d,%d,%v, want %d,%d,%v", tc.placement, start, end, got, tc.start, tc.end, tc.ok)
		}
	}
}

func TestGridAutoFillCountsRepetitions(t *testing.T) {
	got := gridRunIn(t, 800, 200, `
		grid { display: grid; width: 800px; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 10px; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		lay(f, "B")
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// 800px holds three 200px tracks and two 10px gaps; the third track is
	// empty but auto-fill keeps it.
	if got[0].x != 0 || got[0].w != 260 {
		t.Errorf("first = %+v, want 0,0 260 wide", got[0])
	}
	if got[1].x != 270 || got[1].w != 260 {
		t.Errorf("second = %+v, want x 270, 260 wide", got[1])
	}
}

func TestGridAutoFitCollapsesEmptyRepetitions(t *testing.T) {
	got := gridRunIn(t, 800, 200, `
		grid { display: grid; width: 800px; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		lay(f, "B")
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// Two items in a space that fits three tracks: auto-fit drops the empty one
	// and the two tracks share the width.
	if got[0].x != 0 || got[0].w != 395 {
		t.Errorf("first = %+v, want 0,0 395 wide", got[0])
	}
	if got[1].x != 405 || got[1].w != 395 {
		t.Errorf("second = %+v, want x 405, 395 wide", got[1])
	}
}

func TestGridAutoFitCountsAgainstTheSheetScale(t *testing.T) {
	got := gridRunIn(t, 1280, 800, `
		grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		lay(f, "B")
		lay(f, "C")
		lay(f, "D")
	})
	if len(got) != 4 {
		t.Fatalf("got %d boxes, want 4", len(got))
	}
	// The window is drawn at twice the sheet's size, so a 200px minimum is 400
	// on the canvas: the 1280 it has to fit holds three tracks, and the fourth
	// item wraps instead of running off the right edge.
	for i, b := range got {
		if b.x < 0 || b.x+b.w > 1280 {
			t.Errorf("item %d = %+v, want it inside the window", i, b)
		}
	}
	if got[0].y != got[1].y || got[1].y != got[2].y {
		t.Errorf("the first three items = %+v, want them on one row", got[:3])
	}
	if got[3].y <= got[0].y {
		t.Errorf("fourth item = %+v, want it on the row below", got[3])
	}
	if got[3].x != 0 {
		t.Errorf("fourth item = %+v, want it at the left edge", got[3])
	}
}

func TestGridAutoRepeatBetweenFixedTracks(t *testing.T) {
	got := gridRunIn(t, 800, 200, `
		grid { display: grid; width: 800px; grid-template-columns: 100px repeat(auto-fit, minmax(100px, 1fr)); }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		lay(f, "B")
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// 700px is left after the fixed track, enough for seven 100px tracks, and
	// the two items are all auto-fit counts: two repetitions, the second item
	// landing in the first of them.
	if got[0].x != 0 || got[0].w != 100 {
		t.Errorf("first = %+v, want 0,0 100 wide", got[0])
	}
	if got[1].x != 100 || got[1].w != 350 {
		t.Errorf("second = %+v, want 100,0 350 wide", got[1])
	}
}

func TestGridNamedLinesPlaceItems(t *testing.T) {
	got := gridRun(t, `
		grid {
			display: grid; width: 400px;
			grid-template-columns: [full-start] 100px [content-start] 1fr [content-end] 100px [full-end];
			grid-template-rows: [top] 40px [bottom];
		}
		button { grid-column: content-start / content-end; margin: 0; }
		progress { grid-area: top; }
		label { grid-column: -content-end; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "C"); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 3 {
		t.Fatalf("got %d boxes, want 3", len(got))
	}
	// content-start to content-end is the flexible track alone.
	assertGridBox(t, got[0], 100, 0, 200, 40)
	// "top" names a line of the row axis and no line of the column axis, so the
	// progress keeps its column and takes the first.
	assertGridBox(t, got[1], 0, 0, 100, 40)
	// -content-end counts from the end of the axis: the line the last track
	// begins on, with the row left to the flow.
	assertGridBox(t, got[2], 300, 0, 100, got[2].h)
}

func TestGridUnknownLineNameFallsBackToAuto(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; width: 400px; grid-template-columns: 100px 100px; }
		button { grid-column: nowhere; grid-row: nowhere; margin: 0; }
	`, func(f *CSS, dst *[]box) {
		gridLayout(dst)(f, "A")
	})
	if len(got) != 1 {
		t.Fatalf("got %d boxes, want 1", len(got))
	}
	assertGridBox(t, got[0], 0, 0, 100, got[0].h)
}

func TestGridContentKeywordTracks(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; width: 400px; grid-template-columns: min-content max-content fit-content(60px); }
		label { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		gridLayout(dst)(f, "a very long label")
	})
	if len(got) != 1 {
		t.Fatalf("got %d boxes, want 1", len(got))
	}
	// A min-content track holds the widest word, "label", and a max-content
	// track the whole label, 136px of it.
	assertGridBox(t, got[0], 0, 0, 40, got[0].h)
}

func TestGridFitContentClampsBetweenContent(t *testing.T) {
	// fit-content() is the max-content size held between the track's floor and
	// the limit it was given, so a limit below the widest word still leaves the
	// word, and a limit above the label leaves the whole label.
	for _, tc := range []struct {
		limit string
		want  int
	}{
		{"300px", 136},
		{"60px", 60},
		{"40px", 40},
	} {
		got := gridRun(t, `
			grid { display: grid; width: 400px; grid-template-columns: fit-content(`+tc.limit+`); }
			label { margin: 0; }
		`, func(f *CSS, dst *[]box) {
			gridLayout(dst)(f, "a very long label")
		})
		if len(got) != 1 {
			t.Fatalf("limit %s: got %d boxes, want 1", tc.limit, len(got))
		}
		if got[0].w != tc.want {
			t.Errorf("fit-content(%s) = %d, want %d", tc.limit, got[0].w, tc.want)
		}
	}
}

func TestGridAbsoluteItemUsesItsArea(t *testing.T) {
	got := gridRun(t, `
		grid {
			width: 400px; height: 120px; display: grid;
			grid-template-columns: 100px 100px 100px;
			grid-template-rows: 40px 40px 40px;
		}
		progress { position: absolute; grid-area: 2 / 2 / 3 / 3; width: 20px; height: 10px; }
		label { position: absolute; grid-column: 3; grid-row: 1; left: 5px; top: 5px; margin: 0; }
	`, func(f *CSS, dst *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "C"); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// The area is the box the item sits in: with no inset it starts at the
	// area's own corner.
	assertGridBox(t, got[0], 100, 40, 20, 10)
	assertGridBox(t, got[1], 205, 5, got[1].w, got[1].h)
}

func TestGridAbsoluteItemFillsOppositeInsets(t *testing.T) {
	got := gridRun(t, `
		grid {
			width: 400px; height: 120px; display: grid;
			grid-template-columns: 100px 100px 100px;
			grid-template-rows: 40px 40px 40px;
		}
		progress { position: absolute; grid-column: 1 / 3; grid-row: 1; left: 10px; right: 10px; }
	`, func(f *CSS, dst *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 1 {
		t.Fatalf("got %d boxes, want 1", len(got))
	}
	// Two tracks and the gap between them, less both insets, and the height
	// left as the item asked for it since no inset claims the other end.
	assertGridBox(t, got[0], 10, 0, 180, 32)
}

func TestGridItemPercentagesAndAutoSize(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; width: auto; grid-template-columns: 100px 100px; grid-template-rows: 40px; }
		button { width: 50%; }
		progress { width: auto; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 50, 40)
	assertGridBox(t, got["progress"], 100, 0, 100, 40)
}

func TestGridMinimumWinsOverMaximum(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 200px; }
		button { width: 50px; min-width: 100px; max-width: 60px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, got["button"].h)
}

func TestGridRelativeInsets(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px; grid-template-rows: 100px; }
		button { position: relative; left: 5px; top: 7px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 5, 7, got["button"].w, got["button"].h)
}

func TestGridPercentageGapsUseContentBox(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; width: 300px; height: 100px; grid-template-columns: 100px 100px; grid-template-rows: 20px 20px; column-gap: 10%; row-gap: 10%; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["first"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["column"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "C"); ok {
			got["row"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["first"], 50, 0, 100, 20)
	assertGridBox(t, got["column"], 180, 0, 100, 20)
	assertGridBox(t, got["row"], 50, 30, 100, 20)
}

func TestGridFixedTracksAndGap(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; padding: 5px; grid-template-columns: 100px 150px; grid-template-rows: 40px 60px; gap: 10px; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
		lay(f, "C")
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 4 {
		t.Fatalf("got %d boxes, want 4", len(got))
	}
	assertGridBox(t, got[0], 5, 5, 100, 40)
	assertGridBox(t, got[1], 115, 5, 150, 40)
	assertGridBox(t, got[2], 5, 55, 100, 60)
	assertGridBox(t, got[3], 115, 55, 150, 60)
}

func TestGridFlexibleTracksFillWidth(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; grid-template-columns: 1fr 2fr; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		lay(f, "B")
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].w != 133 || got[1].w != 267 {
		t.Errorf("track widths = %d,%d want 133,267", got[0].w, got[1].w)
	}
	assertGridBox(t, got[0], 0, 0, 133, got[0].h)
	assertGridBox(t, got[1], 133, 0, 267, got[1].h)
}

func TestGridTemplateAreasPlaceNamedItems(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 1fr; grid-template-rows: 40px 60px; grid-template-areas: "header header" "side main"; }
		button { grid-area: header; }
		progress { grid-area: side; }
		label { grid-area: main; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "Header"); ok {
			got["header"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["side"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "Main"); ok {
			got["main"] = box{x, y, w, h}
		}
	})
	if len(got) != 3 {
		t.Fatalf("got %d named boxes, want 3", len(got))
	}
	assertGridBox(t, got["header"], 0, 0, 400, 40)
	assertGridBox(t, got["side"], 0, 40, 100, 60)
	assertGridBox(t, got["main"], 100, 40, 300, 60)
}

func TestGridExplicitLinesAndNegativeLines(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: repeat(4, 100px); grid-template-rows: repeat(3, 50px); column-gap: 10px; row-gap: 5px; }
		button { grid-column: 2 / 4; grid-row: 2 / 3; }
		progress { grid-column: 4; grid-row: -2 / -1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	assertGridBox(t, got["button"], 110, 55, 210, 50)
	assertGridBox(t, got["progress"], 330, 110, 100, 50)
}

func TestGridNegativeLineCreatesLeadingImplicitTrack(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 100px; justify-content: start; }
		button { width: 100px; grid-column: -4 / -3; }
		progress { width: 100px; grid-column: 1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, got["button"].h)
	assertGridBox(t, got["progress"], 100, 0, 100, got["progress"].h)
}

func TestGridDensePacksEarlierRow(t *testing.T) {
	sheet := `
		grid { display: grid; grid-template-columns: repeat(3, 100px); grid-template-rows: repeat(2, 50px); grid-auto-flow: row dense; }
		button { grid-column: 2; grid-row: 1; }
		progress { grid-column: 2; }
	`
	got := map[string]box{}
	gridRun(t, sheet, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "Dense"); ok {
			got["label"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["label"], 0, 0, 100, 50)

	sparse := stringsReplaceDense(sheet)
	got = map[string]box{}
	gridRun(t, sparse, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "Sparse"); ok {
			got["label"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["label"], 200, 50, 100, 50)
}

func TestGridAutoPlacementPlacesPartialBeforeAutomatic(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 100px; }
		progress { grid-row: 1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["progress"], 0, 0, got["progress"].w, got["progress"].h)
	assertGridBox(t, got["button"], 100, 0, got["button"].w, got["button"].h)
}

func TestGridPartialPlacementPhasesFollowFlowAxis(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 100px; grid-template-rows: 50px 50px; }
		button { grid-row: 1; }
		progress { grid-column: 1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, 50)
	assertGridBox(t, got["progress"], 0, 50, 100, 50)
}

func TestGridDefiniteColumnKeepsDocumentOrder(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 100px; grid-template-rows: 40px 40px; }
		label { grid-column: 1 / -1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "B"); ok {
			got["label"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, 40)
	assertGridBox(t, got["label"], 0, 40, 200, 40)
}

func TestGridColumnFlowContinuesAfterFixedColumn(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-auto-flow: column; grid-template-columns: 100px 100px; grid-template-rows: repeat(2, 50px); }
		button { grid-column: 1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, 50)
	assertGridBox(t, got["progress"], 0, 50, 100, 50)
}

func TestGridColumnFlowPartialPlacementPhases(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-auto-flow: column; grid-template-columns: 100px 100px; grid-template-rows: 50px 50px; }
		button { grid-column: 1; }
		progress { grid-column: 2; }
		label { grid-row: 1; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "C"); ok {
			got["label"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, 50)
	assertGridBox(t, got["progress"], 100, 50, 100, 50)
	assertGridBox(t, got["label"], 100, 0, 100, 50)
}

func TestGridAutoSpanPlacement(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 100px; }
		button, progress { grid-column: span 2; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 200, got["button"].h)
	assertGridBox(t, got["progress"], 0, got["button"].h, 200, got["progress"].h)
}

func stringsReplaceDense(sheet string) string {
	for i := 0; i+len(" dense") <= len(sheet); i++ {
		if sheet[i:i+len(" dense")] == " dense" {
			return sheet[:i] + sheet[i+len(" dense"):]
		}
	}
	return sheet
}

func TestGridColumnFlowFillsDownThenAcross(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-auto-flow: column; grid-template-columns: repeat(2, 100px); grid-template-rows: repeat(2, 50px); }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "C"); ok {
			got["label"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 100, 50)
	assertGridBox(t, got["progress"], 0, 50, 100, 50)
	assertGridBox(t, got["label"], 100, 0, 100, 50)
}

func TestGridAutoRowsUseIntrinsicHeights(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 50px 50px; }
		button { height: 20px; }
		progress { height: 30px; }
		label { height: 10px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "C"); ok {
			got["label"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, 50, 20)
	assertGridBox(t, got["progress"], 50, 0, 50, 30)
	assertGridBox(t, got["label"], 0, 30, 50, 10)
}

func TestGridFractionRowsNeedDefiniteHeight(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-rows: 1fr 1fr; }
		button { height: 20px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 0, 0, got["button"].w, 20)
}

func TestGridNonStretchContentAlignment(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; width: 300px; grid-template-columns: auto auto; justify-content: end; }
		button { width: 50px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "B"); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 250, 0, 50, got["button"].h)
	assertGridBox(t, got["progress"], 300, 0, 50, got["progress"].h)
}

func TestGridContentAndItemAlignment(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { width: 300px; height: 120px; display: grid; grid-template-columns: 80px 80px; grid-template-rows: 40px 40px; justify-content: end; align-content: end; justify-items: center; align-items: flex-end; }
		button, progress { width: 20px; height: 10px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 220, 70, 20, 10)
	assertGridBox(t, got["progress"], 300, 70, 20, 10)
}

func TestGridAbsoluteItemDoesNotCreateTrack(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { width: 200px; height: 100px; display: grid; }
		progress { position: absolute; left: 10px; top: 20px; width: 30px; height: 15px; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "A"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	assertGridBox(t, got["button"], 100, 0, got["button"].w, got["button"].h)
	assertGridBox(t, got["progress"], 110, 20, 30, 15)
}

func TestGridHiddenItemDoesNotConsumeCell(t *testing.T) {
	got := map[string]box{}
	gridRun(t, `
		grid { display: grid; grid-template-columns: 100px 100px; }
		button { display: none; }
	`, func(f *CSS, _ *[]box) {
		if x, y, w, h, _, ok := f.layout(css.RoleButton, "Hidden"); ok {
			got["button"] = box{x, y, w, h}
		}
		if x, y, w, h, _, ok := f.layout(css.RoleProgress, ""); ok {
			got["progress"] = box{x, y, w, h}
		}
	})
	if _, exists := got["button"]; exists {
		t.Fatal("hidden grid item was painted")
	}
	assertGridBox(t, got["progress"], 0, 0, 100, got["progress"].h)
}

func TestGridContainerFlowsBelow(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "grid { display: grid; grid-template-columns: 100px; }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	gridHeight := 0
	tpl.Grid(func(f *CSS) {
		_, _, _, h, _, ok := f.layout(css.RoleButton, "A")
		if ok {
			gridHeight = h
		}
	})
	_, y, _, _, _, ok := tpl.layout(css.RoleLabel, "after")
	win.End()
	if !ok {
		t.Fatal("label should be visible")
	}
	if gridHeight == 0 || y != gridHeight {
		t.Errorf("label y=%d, grid height=%d", y, gridHeight)
	}
}

func TestGridContainerBoxPaints(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `grid { display: grid; grid-template-columns: 100px; background-color: #ff0000; }`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	win.Begin()
	tpl.Grid(func(f *CSS) { f.layout(css.RoleButton, "A") })
	win.End()
	if got := win.Canvas().At(10, 10); got != canvas.RGBA(255, 0, 0, 255) {
		t.Errorf("container background pixel = %v, want red", got)
	}
}

func TestGridButtonClickInside(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, "grid { display: grid; grid-template-columns: repeat(2, 1fr); }")
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)
	w := textWidth(u, "Go")
	win.Begin()
	testClick(t, win, w/2, textHeight(u)/2)
	var e event.Event
	tpl.Grid(func(f *CSS) {
		e = f.Button("Go")
	})
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("click inside a grid button must produce a Click, got %v", e)
	}
}

func TestGridNestedContainerIsAnItem(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; grid-template-columns: 1fr 1fr; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		// A grid inside the grid is an item of it: the outer cell is its box
		// and its own two columns divide that box in half.
		f.Grid(func(g *CSS) {
			inner := gridLayout(dst)
			inner(g, "B")
			inner(g, "C")
		})
	})
	if len(got) != 3 {
		t.Fatalf("got %d boxes, want 3", len(got))
	}
	assertGridBox(t, got[0], 0, 0, got[0].w, got[0].h)
	assertGridBox(t, got[1], 200, 0, got[1].w, got[1].h)
	assertGridBox(t, got[2], 300, 0, got[2].w, got[2].h)
}

func TestGridNestedContainerSizesItsParentTrack(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; grid-template-columns: auto 1fr; }
		grid { padding: 0; }
		label { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		f.Grid(func(g *CSS) {
			if x, y, w, h, _, ok := g.layout(css.RoleLabel, "a very long label"); ok {
				*dst = append(*dst, box{x, y, w, h})
			}
		})
		if x, y, w, h, _, ok := f.layout(css.RoleLabel, "B"); ok {
			*dst = append(*dst, box{x, y, w, h})
		}
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// The auto column is as wide as the container it holds: the nested grid
	// asks for its label's max-content width, and the flexible column takes
	// what is left.
	if got[0].x != 0 || got[0].w != 136 {
		t.Errorf("nested = %+v, want 0,0 136 wide", got[0])
	}
	if got[1].x != 136 {
		t.Errorf("sibling x = %d, want 136", got[1].x)
	}
}

func TestGridNestedContainerAnswersToItsCell(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; grid-template-columns: 50% 50%; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		lay := gridLayout(dst)
		lay(f, "A")
		f.Grid(func(g *CSS) {
			// A percentage inside the nested grid is half of the box the outer
			// solver handed it, and not half of the 400px it was measured at.
			inner := gridLayout(dst)
			inner(g, "B")
			inner(g, "C")
		})
	})
	if len(got) != 3 {
		t.Fatalf("got %d boxes, want 3", len(got))
	}
	assertGridBox(t, got[0], 0, 0, got[0].w, got[0].h)
	assertGridBox(t, got[1], 200, 0, 100, got[1].h)
	assertGridBox(t, got[2], 300, 0, got[2].w, got[2].h)
}

func TestGridNestedContainerInsideFlex(t *testing.T) {
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		flex { display: flex; }
		grid { display: grid; grid-template-columns: 100px; }
		button { margin: 0; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var got []box
	win.Begin()
	tpl.Flex(func(f *CSS) {
		lay := flexLayout(&got)
		lay(f, "A")
		// A grid inside a flex row is a flex item of its own, taking the
		// 100px its own column asks for.
		f.Grid(func(g *CSS) {
			inner := flexLayout(&got)
			inner(g, "B")
		})
	})
	win.End()
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	if got[0].x != 0 {
		t.Errorf("first x = %d, want 0", got[0].x)
	}
	if got[1].x != got[0].w || got[1].w != 100 {
		t.Errorf("nested = %+v, want x %d 100 wide", got[1], got[0].w)
	}
}

func TestGridNestedButtonClickInside(t *testing.T) {
	win, _, _ := antui.Offscreen(400, 200)
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, `
		grid { display: grid; grid-template-columns: 200px 200px; }
		button { margin: 0; }
	`)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	u := Scale(win)
	w := textWidth(u, "Go")
	win.Begin()
	// The nested grid lands in the second cell, so the button lives at x=200.
	testClick(t, win, 200+w/2, textHeight(u)/2)
	var e event.Event
	tpl.Grid(func(f *CSS) {
		var discard []box
		gridLayout(&discard)(f, "A")
		f.Grid(func(g *CSS) { e = g.Button("Go") })
	})
	win.End()
	if !e.Is(event.Button, event.Click) {
		t.Errorf("click inside a nested grid must produce a Click, got %v", e)
	}
}

func TestGridNestedContainerInsideNestedContainer(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; grid-template-columns: 1fr 1fr; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		f.Grid(func(g *CSS) {
			gridLayout(dst)(g, "A")
			g.Grid(func(h *CSS) {
				inner := gridLayout(dst)
				inner(h, "B")
				inner(h, "C")
			})
		})
	})
	if len(got) != 3 {
		t.Fatalf("got %d boxes, want 3", len(got))
	}
	// The first container fills the page and takes its first cell, so the one
	// inside it is 200px wide and the one inside that is 100px.
	assertGridBox(t, got[0], 0, 0, 100, got[0].h)
	assertGridBox(t, got[1], 100, 0, 50, got[1].h)
	assertGridBox(t, got[2], 150, 0, 50, got[2].h)
}

func TestGridNestedContainerOutOfFlow(t *testing.T) {
	got := gridRun(t, `
		grid { display: grid; grid-template-columns: 1fr 1fr; position: absolute; left: 5px; top: 7px; }
		button { margin: 0; }
	`, func(f *CSS, dst *[]box) {
		gridLayout(dst)(f, "A")
		// The nested container is out of the flow: it takes no cell, and its
		// inset places it against the parent's padding box.
		f.Grid(func(g *CSS) { gridLayout(dst)(g, "B") })
	})
	if len(got) != 2 {
		t.Fatalf("got %d boxes, want 2", len(got))
	}
	// The button keeps its cell; the container inside it is placed by the inset
	// the grid rule gives it.
	assertGridBox(t, got[0], 0, 0, got[0].w, got[0].h)
	assertGridBox(t, got[1], 5, 7, got[1].w, got[1].h)
}
