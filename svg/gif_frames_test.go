package svg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

// An animated GIF in an `<image>` shows the frame the clock calls for, with
// every frame of it read once and kept, instead of the first one forever.

// gifPalette is the colours the frames of these tests are written in, by the
// index each frame paints itself with.
var gifPalette = color.Palette{
	color.NRGBA{R: 255, A: 255},         // 0 red
	color.NRGBA{B: 255, A: 255},         // 1 blue
	color.NRGBA{G: 255, A: 255},         // 2 green
	color.NRGBA{R: 255, G: 255, A: 255}, // 3 yellow
	color.NRGBA{G: 255, B: 255, A: 255}, // 4 cyan
}

// frameOf is a frame painted one colour over its whole rectangle, however
// small the rectangle is: a GIF frame carries only the part it changes, and
// NewPaletted would otherwise leave the index the palette starts at.
func frameOf(r image.Rectangle, idx uint8) *image.Paletted {
	f := image.NewPaletted(r, gifPalette)
	for i := range f.Pix {
		f.Pix[i] = idx
	}
	return f
}

// gifURI is a GIF of the frames given, each shown for the delay beside it in
// hundredths of a second, carried in the href the way a small one is written.
func gifURI(t *testing.T, w, h int, delay []int, disposal []byte, frames ...*image.Paletted) string {
	t.Helper()
	var buf bytes.Buffer
	g := &gif.GIF{
		Image:    frames,
		Delay:    delay,
		Disposal: disposal,
		Config:   image.Config{ColorModel: gifPalette, Width: w, Height: h},
	}
	if err := gif.EncodeAll(&buf, g); err != nil {
		t.Fatalf("encode the gif: %v", err)
	}
	return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

// withClock hands the frame clock to the test: the moment is whatever the
// pointer says, and the real clock goes back when the test is done.
func withClock(t *testing.T, now *float64) {
	t.Helper()
	pictureClock = func() float64 { return *now }
	t.Cleanup(func() { pictureClock = antui.Now })
}

// imageAt is a drawing of one `<image>` filling it, read and checked for
// warnings the way every test here asks for one.
func imageAt(t *testing.T, href string) *Image {
	t.Helper()
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4"/></svg>`, href))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("the drawing said something was wrong: %s", got)
	}
	return img
}

func TestAnAnimatedGIFShowsTheFrameTheClockChooses(t *testing.T) {
	// Two frames of a tenth of a second each, going round a loop of two
	// tenths: the moment on the clock says which one is showing, and the
	// moment before the origin wraps into the loop the same way.
	clock := 0.0
	withClock(t, &clock)
	href := gifURI(t, 4, 4, []int{10, 10}, nil, frameOf(image.Rect(0, 0, 4, 4), 0), frameOf(image.Rect(0, 0, 4, 4), 1))
	img := imageAt(t, href)
	for _, tc := range []struct {
		name string
		at   float64
		want canvas.Color
	}{
		{"on the first frame", 0, canvas.RGB(255, 0, 0)},
		{"still on it", 0.05, canvas.RGB(255, 0, 0)},
		{"the second frame", 0.15, canvas.RGB(0, 0, 255)},
		{"round the loop", 0.25, canvas.RGB(255, 0, 0)},
		{"wrapped from before the origin", -0.05, canvas.RGB(0, 0, 255)},
	} {
		clock = tc.at
		wantColour(t, img.Render(4, 4), 2, 2, tc.want)
	}
}

func TestAnAnimatedGIFKeepsItsFramesInTheCache(t *testing.T) {
	// The frames are read once, under the address they came from, and the
	// canvas of the moment is kept while the clock still shows that moment:
	// asking again hands the same one back, and the clock moving on is what
	// makes a new paint rather than the frame that has ended.
	clock := 0.05
	withClock(t, &clock)
	href := gifURI(t, 4, 4, []int{10, 10}, nil, frameOf(image.Rect(0, 0, 4, 4), 0), frameOf(image.Rect(0, 0, 4, 4), 1))
	img := imageAt(t, href)
	if len(img.pictures) != 1 {
		t.Fatalf("the picture was read %d times, want one read under its address", len(img.pictures))
	}
	for _, f := range img.pictures {
		if len(f.frames) != 2 {
			t.Errorf("the picture was read into %d frames, want the 2 the GIF has", len(f.frames))
		}
	}
	here := img.Render(4, 4)
	if again := img.Render(4, 4); again != here {
		t.Error("the canvas of the moment was made again while the clock had not moved")
	}
	clock = 0.15
	moved := img.Render(4, 4)
	if moved == here {
		t.Error("the canvas of the frame that ended came back after the clock moved on")
	}
	wantColour(t, moved, 2, 2, canvas.RGB(0, 0, 255))
}

func TestAnAnimatedGIFDrawsEachFrameOnWhatTheLastOneLeft(t *testing.T) {
	// A GIF stores the change each frame makes, not the whole picture, and the
	// disposal byte says what of it is left for the next one: five frames
	// going over, restoring what was there, and clearing what they covered.
	clock := 0.05
	withClock(t, &clock)
	href := gifURI(t, 4, 4,
		[]int{10, 10, 10, 10, 10},
		[]byte{gif.DisposalNone, gif.DisposalPrevious, gif.DisposalNone, gif.DisposalBackground, gif.DisposalNone},
		frameOf(image.Rect(0, 0, 4, 4), 0), // red, the whole square
		frameOf(image.Rect(1, 1, 3, 3), 1), // blue on top, and gone after
		frameOf(image.Rect(0, 0, 2, 2), 2), // green over what was restored
		frameOf(image.Rect(2, 2, 4, 4), 3), // yellow, cleared after
		frameOf(image.Rect(3, 3, 4, 4), 4), // cyan where the yellow was
	)
	img := imageAt(t, href)
	for _, tc := range []struct {
		name string
		at   float64
		x, y int
		want canvas.Color
	}{
		{"the first frame is the whole square", 0.05, 3, 3, canvas.RGB(255, 0, 0)},
		{"the change is drawn on top", 0.15, 2, 2, canvas.RGB(0, 0, 255)},
		{"the restored frame starts from the whole square", 0.25, 2, 2, canvas.RGB(255, 0, 0)},
		{"the change beside it stays", 0.25, 0, 0, canvas.RGB(0, 255, 0)},
		{"the frame about to clear draws its own change", 0.35, 3, 3, canvas.RGB(255, 255, 0)},
		{"what the clear left is nothing where the yellow was", 0.45, 2, 2, canvas.Transparent},
		{"and draws its own change", 0.45, 3, 3, canvas.RGB(0, 255, 255)},
	} {
		clock = tc.at
		wantColour(t, img.Render(4, 4), tc.x, tc.y, tc.want)
	}
}
