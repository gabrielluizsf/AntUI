package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// An icon that is a moving GIF goes back to the drawing as the clock runs:
// the canvas of the frame that ended is not worth remembering, so IconCanvas
// keeps nothing of it and the next frame is the next paint.

// movingGIFURI is a two-frame GIF of eight pixels a side, red and blue, each
// frame shown for a tenth of a second — a loop short enough for a test to
// watch it turn without timing where on it the clock happens to be — carried
// in the href of the `<image>` that draws it.
func movingGIFURI(t *testing.T) string {
	t.Helper()
	pal := color.Palette{color.NRGBA{R: 255, A: 255}, color.NRGBA{B: 255, A: 255}}
	frames := make([]*image.Paletted, 2)
	for i := range frames {
		f := image.NewPaletted(image.Rect(0, 0, 8, 8), pal)
		for j := range f.Pix {
			f.Pix[j] = uint8(i)
		}
		frames[i] = f
	}
	var buf bytes.Buffer
	err := gif.EncodeAll(&buf, &gif.GIF{
		Image:  frames,
		Delay:  []int{10, 10},
		Config: image.Config{ColorModel: pal, Width: 8, Height: 8},
	})
	if err != nil {
		t.Fatalf("encode the gif: %v", err)
	}
	return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

// TestAnAnimatedIconKeepsMoving asks for the icon over and over the way a frame
// does and requires both of its frames to come round: the clock may start the
// test anywhere on the loop, so what proves the movement is seeing the two
// colours rather than either one in particular. One colour alone, for as long
// as the test is willing to wait, is the picture of the first frame kept for
// ever — which is what this is here to say does not happen.
func TestAnAnimatedIconKeepsMoving(t *testing.T) {
	restoreIcons(t)
	src := fmt.Sprintf(
		`<svg viewBox="0 0 8 8"><image href="%s" x="0" y="0" width="8" height="8"/></svg>`,
		movingGIFURI(t))
	if err := SetIcon("gifmoving", src); err != nil {
		t.Fatalf("register the icon: %v", err)
	}
	area := canvas.Area{Width: 8, Height: 8}
	seen := map[canvas.Color]bool{}
	deadline := time.Now().Add(5 * time.Second)
	for {
		cv := IconCanvas("gifmoving", area, canvas.RGB(255, 0, 0))
		if cv == nil {
			t.Fatal("the icon was not painted at all")
		}
		seen[cv.At(4, 4)] = true
		if len(seen) == 2 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("the icon never left one frame: pixel (4,4) stayed %v", cv.At(4, 4))
		}
		time.Sleep(time.Millisecond)
	}
}
