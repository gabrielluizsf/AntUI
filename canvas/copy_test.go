package canvas

import "testing"

// opaqueTile is a picture with no two pixels alike and no alpha to blend.
func opaqueTile(t testing.TB, w, h int) *Canvas {
	t.Helper()
	tile, err := NewCanvas(w, h)
	if err != nil {
		t.Fatalf("new canvas: %v", err)
	}
	for i := range tile.Pixels {
		tile.Pixels[i] = Color(0xFF000000 | uint32(i*2654435761)&0xFFFFFF)
	}
	return tile
}

// Over an opaque picture Copy is Blit: same pixels, at every placement that
// crosses an edge or the clip, on every format.
func TestCopyOfAnOpaquePictureIsBlit(t *testing.T) {
	tile := opaqueTile(t, 13, 9)
	for _, format := range []Format{ARGB32, RGB565, Pal8} {
		for _, at := range []Point{{X: 0, Y: 0}, {X: 5, Y: 3}, {X: -4, Y: -2}, {X: 27, Y: 18}, {X: 40, Y: 40}} {
			blit, _ := NewCanvasFormat(32, 24, format)
			copied, _ := NewCanvasFormat(32, 24, format)
			for _, cv := range []*Canvas{blit, copied} {
				cv.Clear(RGB(0x10, 0x20, 0x30))
				cv.SetClip(Area{X: 2, Y: 1, Width: 26, Height: 20})
			}
			blit.Blit(int(at.X), int(at.Y), tile)
			copied.Copy(int(at.X), int(at.Y), tile)
			for y := range 24 {
				for x := range 32 {
					if got, want := copied.At(x, y), blit.At(x, y); got != want {
						t.Fatalf("format %d at %v: pixel (%d,%d) = %08X, want %08X",
							format, at, x, y, uint32(got), uint32(want))
					}
				}
			}
		}
	}
}

// Copy replaces rather than blends: a translucent pixel lands as it is.
func TestCopyDoesNotBlend(t *testing.T) {
	cv, _ := NewCanvas(4, 4)
	cv.Clear(White)
	src, _ := NewCanvas(2, 2)
	half := RGBA(RGBAParams{R: 0xFF, A: 0x80})
	for i := range src.Pixels {
		src.Pixels[i] = half
	}
	cv.Copy(1, 1, src)
	if got := cv.At(1, 1); got != half {
		t.Errorf("copied pixel = %08X, want %08X", uint32(got), uint32(half))
	}
	if got := cv.At(0, 0); got != White {
		t.Errorf("pixel outside the copy = %08X, want white", uint32(got))
	}
	cv.Copy(0, 0, nil)
}

func BenchmarkCopyVsBlit(b *testing.B) {
	tile := opaqueTile(b, 64, 64)
	dst, _ := NewCanvas(1920, 1080)
	screen := func(draw func(x, y int, src *Canvas)) {
		for y := 0; y < 1080; y += 64 {
			for x := 0; x < 1920; x += 64 {
				draw(x, y, tile)
			}
		}
	}
	b.Run("blit", func(b *testing.B) {
		for b.Loop() {
			screen(dst.Blit)
		}
	})
	b.Run("copy", func(b *testing.B) {
		for b.Loop() {
			screen(dst.Copy)
		}
	})
}

// A clip set by hand wider than the canvas must not let a row run on into
// the next one: the write stops at the canvas's own edge.
func TestCopyStaysInsideTheRow(t *testing.T) {
	cv, _ := NewCanvas(4, 2)
	cv.Clip = Area{Width: 8, Height: 2}
	defer func() {
		recover()
		for x := range 4 {
			if got := cv.At(x, 1); got != 0 {
				t.Fatalf("pixel (%d,1) = %08X: row 0 ran into row 1", x, uint32(got))
			}
		}
	}()
	cv.Copy(0, 0, opaqueTile(t, 8, 1))
}
