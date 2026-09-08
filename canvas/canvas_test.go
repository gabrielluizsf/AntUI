package canvas

import "testing"

func newCanvas(t *testing.T, w, h int) *Canvas {
	t.Helper()
	cv, err := NewCanvas(w, h)
	if err != nil {
		t.Fatalf("NewCanvas(%d, %d): %v", w, h, err)
	}
	return cv
}

func TestCanvasBasics(t *testing.T) {
	cv := newCanvas(t, 16, 8)
	if cv.Width != 16 || cv.Height != 8 {
		t.Fatalf("size = %dx%d, want 16x8", cv.Width, cv.Height)
	}

	cv.Clear(Red)
	if got := cv.At(0, 0); got != Red {
		t.Errorf("first pixel = %08X, want the clear colour", uint32(got))
	}
	if got := cv.At(15, 7); got != Red {
		t.Errorf("last pixel = %08X, want the clear colour", uint32(got))
	}

	cv.Pixel(3, 4, Blue)
	if got := cv.At(3, 4); got != Blue {
		t.Errorf("painted pixel = %08X, want blue", uint32(got))
	}

	cv.FillRect(0, 0, 4, 4, Green)
	if got := cv.At(3, 3); got != Green {
		t.Errorf("inside the rectangle = %08X, want green", uint32(got))
	}
	if got := cv.At(4, 4); got != Red {
		t.Errorf("outside the rectangle = %08X, want it untouched", uint32(got))
	}
}

func TestNewCanvasRejectsBadSize(t *testing.T) {
	for _, size := range [][2]int{{0, 10}, {10, 0}, {-1, 5}} {
		if _, err := NewCanvas(size[0], size[1]); err == nil {
			t.Errorf("NewCanvas(%d, %d) returned no error", size[0], size[1])
		}
	}
}

func TestDrawingOutOfBoundsIsClipped(t *testing.T) {
	cv := newCanvas(t, 8, 8)
	cv.Clear(Black)

	// None of these may panic, and none may paint outside the canvas.
	cv.Pixel(-5, -5, White)
	cv.Pixel(100, 100, White)
	cv.FillRect(-20, -20, 5, 5, White)
	cv.Line(-50, -50, 60, 60, White)
	cv.Text(-100, -100, "outside", White)

	if got := cv.At(0, 0); got != White {
		t.Errorf("origin = %08X, want the diagonal line to cross it", uint32(got))
	}
	if got := cv.At(7, 0); got != Black {
		t.Errorf("corner = %08X, want it untouched", uint32(got))
	}
}

func TestClipping(t *testing.T) {
	cv := newCanvas(t, 10, 10)
	cv.Clear(Black)
	cv.SetClip(2, 2, 4, 4)
	cv.FillRect(0, 0, 10, 10, White)

	if got := cv.At(1, 1); got != Black {
		t.Errorf("outside the clip = %08X, want it untouched", uint32(got))
	}
	if got := cv.At(3, 3); got != White {
		t.Errorf("inside the clip = %08X, want white", uint32(got))
	}
	if got := cv.At(6, 6); got != Black {
		t.Errorf("past the clip = %08X, want it untouched", uint32(got))
	}

	cv.ResetClip()
	cv.FillRect(0, 0, 10, 10, White)
	if got := cv.At(9, 9); got != White {
		t.Errorf("after ResetClip = %08X, want everything painted", uint32(got))
	}
}

func TestClipIsIntersectedWithTheCanvas(t *testing.T) {
	cv := newCanvas(t, 10, 10)
	cv.SetClip(-5, -5, 100, 100)
	if cv.Clip != (Area{0, 0, 10, 10}) {
		t.Errorf("clip = %+v, want it trimmed to the canvas", cv.Clip)
	}
	cv.SetClip(20, 20, 5, 5)
	if cv.Clip.Width != 0 || cv.Clip.Height != 0 {
		t.Errorf("clip = %+v, want it empty when it falls off the canvas", cv.Clip)
	}
}

func TestText(t *testing.T) {
	cv := newCanvas(t, 64, 20)
	cv.Clear(Black)

	if got := cv.Text(0, 0, "AB", White); got != 16 {
		t.Errorf("Text width = %d, want 16", got)
	}
	lit := 0
	for y := range 16 {
		for x := range 16 {
			if cv.At(x, y) == White {
				lit++
			}
		}
	}
	if lit <= 10 {
		t.Errorf("lit pixels = %d, want the letters to leave a mark", lit)
	}

	cv.Clear(Black)
	cv.Text(0, 0, " ", White)
	if got := cv.At(3, 8); got != Black {
		t.Errorf("a space lit pixel (%08X), want none", uint32(got))
	}

	if got := cv.TextScaled(0, 0, "A", White, 2); got != 16 {
		t.Errorf("TextScaled(2) width = %d, want 16", got)
	}
}

func TestTextWidth(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"", 0},
		{"A", FontWidth},
		{"AB", 2 * FontWidth},
		{"\t", 4 * FontWidth},
		{"ação", 4 * FontWidth},     // accents are one cell each
		{"AB\nABCD", 4 * FontWidth}, // the widest line wins
		{"ABCD\nAB", 4 * FontWidth}, // whichever side it is on
	}
	for _, tt := range tests {
		if got := TextWidth(tt.text); got != tt.want {
			t.Errorf("TextWidth(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}

func TestBlit(t *testing.T) {
	dst := newCanvas(t, 10, 10)
	src := newCanvas(t, 4, 4)
	dst.Clear(Black)
	src.Clear(Red)
	dst.Blit(3, 3, src)

	if got := dst.At(3, 3); got != Red {
		t.Errorf("start of the copy = %08X, want red", uint32(got))
	}
	if got := dst.At(6, 6); got != Red {
		t.Errorf("end of the copy = %08X, want red", uint32(got))
	}
	if got := dst.At(7, 7); got != Black {
		t.Errorf("past the copy = %08X, want it untouched", uint32(got))
	}
}

func TestShapes(t *testing.T) {
	cv := newCanvas(t, 40, 40)

	t.Run("filled circle", func(t *testing.T) {
		cv.Clear(Black)
		cv.FillCircle(20, 20, 10, White)
		if got := cv.At(20, 20); got != White {
			t.Errorf("centre = %08X, want it filled", uint32(got))
		}
		if got := cv.At(0, 0); got != Black {
			t.Errorf("corner = %08X, want it outside the circle", uint32(got))
		}
	})

	t.Run("filled triangle", func(t *testing.T) {
		cv.Clear(Black)
		cv.FillTriangle(0, 0, 39, 0, 0, 39, White)
		if got := cv.At(2, 2); got != White {
			t.Errorf("inside = %08X, want it painted", uint32(got))
		}
		if got := cv.At(38, 38); got != Black {
			t.Errorf("outside = %08X, want it clean", uint32(got))
		}
	})

	t.Run("filled rounded rectangle", func(t *testing.T) {
		cv.Clear(Black)
		cv.FillRoundRect(0, 0, 40, 40, 8, White)
		if got := cv.At(20, 20); got != White {
			t.Errorf("middle = %08X, want it filled", uint32(got))
		}
		if got := cv.At(0, 0); got != Black {
			t.Errorf("corner = %08X, want it cut away", uint32(got))
		}
	})
}

// The narrow formats have to survive everything the wide one does. The C
// original crashed here: its circle and rounded-rect fills wrote through the
// 32-bit pixel pointer, which is NULL on a narrow canvas.
func TestNarrowFormats(t *testing.T) {
	for _, format := range []Format{RGB565, Pal8} {
		t.Run(format.String(), func(t *testing.T) {
			cv, err := NewCanvasFormat(32, 32, format)
			if err != nil {
				t.Fatalf("NewCanvasFormat: %v", err)
			}
			if cv.Pixels != nil {
				t.Error("Pixels should be nil on a narrow canvas, so that code reaching into it draws nothing rather than corrupting the buffer")
			}
			cv.Clear(Black)
			cv.FillRect(0, 0, 32, 32, Black)
			cv.FillCircle(16, 16, 8, White)
			cv.FillRoundRect(2, 2, 12, 12, 4, White)
			cv.Text(0, 0, "hi", White)

			if got := cv.At(16, 16); got.R() < 200 || got.G() < 200 || got.B() < 200 {
				t.Errorf("circle centre = %08X, want it near white", uint32(got))
			}
			if got := cv.At(31, 31); got.R() > 60 {
				t.Errorf("far corner = %08X, want it near black", uint32(got))
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	for format, want := range map[Format]int{ARGB32: 4, RGB565: 2, Pal8: 1} {
		if got := format.Bytes(); got != want {
			t.Errorf("%s.Bytes() = %d, want %d", format, got, want)
		}
	}
}

func TestPalette(t *testing.T) {
	cv, err := NewCanvasFormat(8, 8, Pal8)
	if err != nil {
		t.Fatalf("NewCanvasFormat: %v", err)
	}
	if err := cv.SetPalette([]Color{Red, Green, Blue}); err != nil {
		t.Fatalf("SetPalette: %v", err)
	}
	if got := cv.PaletteIndex(Red); got != 0 {
		t.Errorf("PaletteIndex(Red) = %d, want 0", got)
	}
	// Fewer than 256 entries are padded with the last one, so anything past
	// the third index is blue.
	if got := cv.PaletteIndex(Blue); got < 2 {
		t.Errorf("PaletteIndex(Blue) = %d, want the blue entry", got)
	}

	cv.Clear(Green)
	if got := cv.At(4, 4); got != Green|0xFF000000 {
		t.Errorf("cleared pixel = %08X, want the palette's green", uint32(got))
	}
}

func TestSetPaletteRejectsWideCanvas(t *testing.T) {
	cv := newCanvas(t, 8, 8)
	if err := cv.SetPalette([]Color{Red}); err == nil {
		t.Error("SetPalette on an ARGB32 canvas returned no error")
	}
}

func TestView(t *testing.T) {
	sheet := newCanvas(t, 16, 16)
	sheet.Clear(Black)

	view := sheet.View(4, 4, 8, 8)
	if view.Width != 8 || view.Height != 8 {
		t.Fatalf("view is %dx%d, want 8x8", view.Width, view.Height)
	}
	if view.Stride != sheet.Stride {
		t.Errorf("view stride = %d, want the parent's %d", view.Stride, sheet.Stride)
	}

	// Writing through the view writes through to the sheet, at the offset.
	view.Clear(Red)
	if got := sheet.At(4, 4); got != Red {
		t.Errorf("sheet at the view's origin = %08X, want red", uint32(got))
	}
	if got := sheet.At(11, 11); got != Red {
		t.Errorf("sheet at the view's far corner = %08X, want red", uint32(got))
	}
	// And nothing outside it.
	if got := sheet.At(3, 3); got != Black {
		t.Errorf("sheet outside the view = %08X, want it untouched", uint32(got))
	}
	if got := sheet.At(12, 12); got != Black {
		t.Errorf("sheet past the view = %08X, want it untouched", uint32(got))
	}
}

func TestViewIsClampedToTheParent(t *testing.T) {
	sheet := newCanvas(t, 16, 16)

	view := sheet.View(10, 10, 100, 100)
	if view.Width != 6 || view.Height != 6 {
		t.Errorf("view is %dx%d, want it clamped to 6x6", view.Width, view.Height)
	}
	// Wholly outside gives an empty view rather than a panic.
	if got := sheet.View(50, 50, 4, 4); got.Width != 0 || got.Height != 0 {
		t.Errorf("a view off the sheet is %dx%d, want it empty", got.Width, got.Height)
	}
	// Drawing into an empty view must not panic either.
	sheet.View(50, 50, 4, 4).Clear(Red)
}

func TestNarrowCanvasHasNoView(t *testing.T) {
	cv, err := NewCanvasFormat(16, 16, RGB565)
	if err != nil {
		t.Fatalf("NewCanvasFormat: %v", err)
	}
	if got := cv.View(0, 0, 8, 8); got != nil {
		t.Error("a narrow canvas has no view to give; nil says so, an empty canvas would not")
	}
}
