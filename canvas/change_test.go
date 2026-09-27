package canvas

import "testing"

// TestDamageCoversEveryDrawing checks the promise the baseline makes: whatever
// a program draws, the rectangle PresentChanges answers covers every pixel
// that really moved. A drawing path that forgets to say it wrote a row would
// leave that row off the damage and the screen would keep showing the frame
// before it, so every path in the package is walked here.
func TestDamageCoversEveryDrawing(t *testing.T) {
	cv, err := NewCanvas(320, 240)
	if err != nil {
		t.Fatalf("new canvas: %v", err)
	}
	paints := map[string]func(){
		"clear":        func() { cv.Clear(RGB(0x10, 0x12, 0x18)) },
		"fill rect":    func() { cv.FillRect(20, 20, 200, 120, RGB(0x3E, 0x63, 0xDD)) },
		"rect":         func() { cv.Rect(4, 4, 300, 220, RGBA(0xFF, 0xFF, 0xFF, 0x80)) },
		"line":         func() { cv.Line(0, 239, 319, 0, Red) },
		"circle":       func() { cv.Circle(160, 120, 60, White) },
		"fill circle":  func() { cv.FillCircle(80, 80, 40, Green) },
		"fill ellipse": func() { cv.FillEllipse(240, 80, 40, 20, Blue) },
		"round rect": func() {
			cv.FillRoundRect(20, 150, 120, 60, 12, Yellow)
			cv.RoundRect(160, 150, 120, 60, 12, Black)
		},
		"triangle": func() { cv.FillTriangle(10, 220, 60, 220, 35, 180, Cyan) },
		"ring arc": func() { cv.RingArc(160, 120, 70, 8, 0.3, 2.0, RGB(0xFF, 0x88, 0x00)) },
		"text":     func() { cv.Text(20, 30, "AntUI 24:00", White) },
		"face":     func() { cv.face().Draw(cv, 20, 50, "AntUI", White) },
		"styled":   func() { cv.DrawStyled(20, 70, "styled", White, TextStyle{Bold: true, Scale: 2}) },
		"gradient": func() {
			cv.FillGradient(220, 20, 80, 60, Gradient{Kind: GLinear, Stops: []GradientStop{
				{Offset: 0, Color: Red}, {Offset: 1, Color: Blue}}})
		},
		"mask rect": func() {
			cv.FillRect(0, 100, 320, 40, RGB(0x22, 0x44, 0x66))
			cv.MaskRect(10, 100, 100, 40)
		},
		"mask round rect": func() {
			cv.FillRect(200, 100, 100, 40, RGB(0x66, 0x44, 0x22))
			cv.MaskRoundRect(200, 100, 100, 40, 16, 16)
		},
		"blur":        func() { cv.Blur(10, 10, 80, 60, 3) },
		"filter":      func() { cv.FilterRegion(100, 10, 60, 40, FilterBlur, 2) },
		"scale alpha": func() { ScaleAlpha(cv, 0.5) },
		"box shadow":  func() { cv.BoxShadow(20, 20, 100, 50, 8, 8, 4, 4, 6, 2, RGBA(0, 0, 0, 0x99), false) },
		"blit":        func() { cv.Blit(120, 30, layer(t)) },
		"blit over":   func() { cv.BlitOver(140, 30, layer(t)) },
		"blit scaled": func() { cv.BlitScaled(160, 30, 64, 64, layer(t)) },
		"blit matrix": func() {
			cv.BlitMatrix(layer(t), Matrix{A: 1, D: 1}, Area{Width: 32, Height: 32})
		},
		"backdrop": func() {
			cv.BlitBackdrop(layer(t), Matrix{A: 1, D: 1}, Area{X: 0, Y: 100, Width: 200, Height: 40})
		},
		"put": func() { cv.Put(300, 10, RGB(1, 2, 3)) },
	}

	base := make([]Color, cv.Width*cv.Height)
	for name, paint := range paints {
		// A frame of its own for every paint, so one cannot hide another.
		cv.ResetClip()
		cv.Clear(RGB(0x10, 0x12, 0x18))
		cv.FillRect(20, 20, 200, 120, RGB(0x3E, 0x63, 0xDD))
		copy(base, cv.Pixels)
		cv.Compare(base)

		paint()
		dirty := cv.PresentChanges()

		missed, first := 0, [2]int{}
		for y := range cv.Height {
			for x := range cv.Width {
				at := y*cv.Width + x
				if cv.Pixels[y*cv.Stride+x] == base[at] {
					continue
				}
				if x >= dirty.X && x < dirty.X+dirty.Width &&
					y >= dirty.Y && y < dirty.Y+dirty.Height {
					continue
				}
				if missed == 0 {
					first = [2]int{x, y}
				}
				missed++
			}
		}
		if missed > 0 {
			t.Errorf("%s: %d pixels changed outside the damage %v, first at (%d,%d)",
				name, missed, dirty, first[0], first[1])
		}
		if cv.Comparing() != true {
			t.Errorf("%s: the baseline was dropped", name)
		}
	}
}

// TestDamageIgnoresPaintThatWasPaintedOver is the other half: drawing is
// overdraw, and a frame that ends up looking like the one before it has
// nothing to send however much was painted on the way there.
func TestDamageIgnoresPaintThatWasPaintedOver(t *testing.T) {
	cv, err := NewCanvas(64, 64)
	if err != nil {
		t.Fatalf("new canvas: %v", err)
	}
	cv.Clear(Black)
	base := make([]Color, cv.Width*cv.Height)
	copy(base, cv.Pixels)
	cv.Compare(base)

	// A shadow, then the card that covers it, then a border on top: three
	// writes over the same pixels that finish where they started.
	overdraw := func() {
		cv.FillRect(10, 10, 40, 40, RGBA(0, 0, 0, 0x40))
		cv.FillRect(12, 12, 36, 36, White)
		cv.Rect(12, 12, 36, 36, Black)
	}
	overdraw()
	if dirty := cv.PresentChanges(); dirty.Width == 0 {
		t.Error("the first pass changed the frame and should have been sent")
	}
	overdraw()
	overdraw()
	if dirty := cv.PresentChanges(); dirty.Width != 0 {
		t.Errorf("damage = %v, want nothing: the frame ended where it began", dirty)
	}

	// And one pixel that really moves is found, with nothing around it.
	cv.Put(30, 30, Red)
	dirty := cv.PresentChanges()
	if dirty != (Area{X: 30, Y: 30, Width: 1, Height: 1}) {
		t.Errorf("damage = %v, want the single pixel at 30,30", dirty)
	}
}

// TestABaselineOfTheWrongSizeTurnsTheMeasuringOff is what a resized canvas
// gets: there is nothing to compare against, so everything it draws is new.
func TestABaselineOfTheWrongSizeTurnsTheMeasuringOff(t *testing.T) {
	cv, err := NewCanvas(32, 32)
	if err != nil {
		t.Fatalf("new canvas: %v", err)
	}
	cv.Compare(make([]Color, 16*16))
	if cv.Comparing() {
		t.Error("a baseline that does not fit was taken")
	}
	cv.Clear(White)
	if dirty := cv.PresentChanges(); dirty.Width != 0 {
		t.Errorf("damage = %v, want nothing without a baseline", dirty)
	}
	if cv.Dirty.Width == 0 {
		t.Error("the dirty bounds should still say what was written")
	}
}

func layer(t *testing.T) *Canvas {
	t.Helper()
	l, err := NewCanvas(32, 32)
	if err != nil {
		t.Fatalf("new layer: %v", err)
	}
	l.Clear(RGB(0x33, 0x99, 0x66))
	l.Circle(16, 16, 12, Yellow)
	return l
}
