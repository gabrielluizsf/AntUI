package canvas

import "testing"

func inkCount(cv *Canvas) int {
	n := 0
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			if cv.At(x, y).A() != 0 {
				n++
			}
		}
	}
	return n
}

// topInkRow returns the highest row that holds any ink, or -1.
func topInkRow(cv *Canvas) int {
	for y := 0; y < cv.Height; y++ {
		for x := 0; x < cv.Width; x++ {
			if cv.At(x, y).A() != 0 {
				return y
			}
		}
	}
	return -1
}

// leftmostAt returns the first painted column of a row, or -1.
func leftmostAt(cv *Canvas, y int) int {
	for x := 0; x < cv.Width; x++ {
		if cv.At(x, y).A() != 0 {
			return x
		}
	}
	return -1
}

func TestTextWidthStyledSpacing(t *testing.T) {
	base := TextWidthStyled("a b", TextStyle{})
	if got := TextWidthStyled("a b", TextStyle{LetterSpacing: 3}); got <= base {
		t.Fatalf("letter-spacing should widen, got %d base %d", got, base)
	}
	if got := TextWidthStyled("a b", TextStyle{WordSpacing: 5}); got <= base {
		t.Fatalf("word-spacing should widen, got %d base %d", got, base)
	}
}

func TestDrawStyledWidthMatches(t *testing.T) {
	cv, _ := NewCanvas(200, 40)
	o := TextStyle{LetterSpacing: 2, WordSpacing: 4}
	want := cv.TextWidthStyled("hello world", o)
	if got := cv.DrawStyled(0, 30, "hello world", White, o); got != want {
		t.Fatalf("draw width %d, measured %d", got, want)
	}
}

func TestDrawStyledBoldAddsInk(t *testing.T) {
	plain, _ := NewCanvas(60, 40)
	bold, _ := NewCanvas(60, 40)
	plain.DrawStyled(0, 30, "H", White, TextStyle{})
	bold.DrawStyled(0, 30, "H", White, TextStyle{Bold: true})
	if inkCount(bold) <= inkCount(plain) {
		t.Fatalf("bold should add ink, plain %d bold %d", inkCount(plain), inkCount(bold))
	}
}

func TestDrawStyledItalicLeans(t *testing.T) {
	plain, _ := NewCanvas(60, 40)
	italic, _ := NewCanvas(60, 40)
	plain.DrawStyled(10, 10, "H", White, TextStyle{})
	italic.DrawStyled(10, 10, "H", White, TextStyle{Italic: true})
	py, iy := topInkRow(plain), topInkRow(italic)
	if py < 0 || iy < 0 {
		t.Fatal("expected some ink from either draw")
	}
	if leftmostAt(italic, iy) <= leftmostAt(plain, py) {
		t.Fatalf("italic top row should lean right, plain %d italic %d",
			leftmostAt(plain, py), leftmostAt(italic, iy))
	}
}
