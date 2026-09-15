// The builder is host-side tooling. Its tests are platform-independent on
// Linux and macOS, but on Windows a git checkout turns the golden files
// under testdata into CRLF, so byte comparisons against them fail. These
// are Android tests; let the Android workflow run them.
//
//go:build !windows

package apk

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// source writes a square picture with an opaque middle and transparent
// corners, which is what a launcher icon usually is.
func source(t *testing.T, dir string, size int) string {
	t.Helper()
	cv, err := canvas.NewCanvas(size, size)
	if err != nil {
		t.Fatal(err)
	}
	cv.Clear(canvas.Transparent)
	cv.FillCircle(size/2, size/2, size/2, canvas.Blue)
	path := filepath.Join(dir, "src.png")
	if err := writePNG(path, cv); err != nil {
		t.Fatal(err)
	}
	return path
}

func read(t *testing.T, path string) (w, h int, alphaAt func(x, y int) uint8) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	b := img.Bounds()
	return b.Dx(), b.Dy(), func(x, y int) uint8 {
		_, _, _, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
		return uint8(a >> 8)
	}
}

// A round icon has to actually be round: a launcher that asks for one is a
// launcher that will not mask what it is given, so the corners have to be
// gone before it gets there.
func TestRoundIconHasNoCorners(t *testing.T) {
	dir := t.TempDir()
	if err := WriteIcons(dir, source(t, dir, 256), canvas.Blue); err != nil {
		t.Fatal(err)
	}
	w, h, alpha := read(t, filepath.Join(dir, "mipmap-xxxhdpi", "ic_launcher_round.png"))
	if w != 192 || h != 192 {
		t.Fatalf("the xxxhdpi icon is %dx%d, want 192x192", w, h)
	}
	for _, c := range [][2]int{{1, 1}, {w - 2, 1}, {1, h - 2}, {w - 2, h - 2}} {
		if a := alpha(c[0], c[1]); a != 0 {
			t.Errorf("the corner at %v has alpha %d, and a round icon has no corners", c, a)
		}
	}
	if a := alpha(w/2, h/2); a != 255 {
		t.Errorf("the middle has alpha %d, want 255", a)
	}
}

// The adaptive foreground is a 108dp layer with only the middle 72dp certain
// to be seen. Anything drawn in the outer sixth may be cropped by the
// launcher, so nothing should be there.
func TestAdaptiveLayerKeepsItsMargin(t *testing.T) {
	dir := t.TempDir()
	if err := WriteIcons(dir, source(t, dir, 256), canvas.Blue); err != nil {
		t.Fatal(err)
	}
	w, h, alpha := read(t, filepath.Join(dir, "mipmap-xxxhdpi", "ic_launcher_foreground.png"))
	if w != 432 || h != 432 {
		t.Fatalf("the xxxhdpi foreground is %dx%d, want 432x432", w, h)
	}
	margin := int(float64(w) * (1 - safeFraction) / 2)
	for x := 0; x < w; x += 8 {
		if a := alpha(x, margin/2); a != 0 {
			t.Fatalf("the layer has alpha %d at (%d, %d), inside the %d-pixel "+
				"margin a launcher may crop", a, x, margin/2, margin)
		}
	}
	if a := alpha(w/2, h/2); a != 255 {
		t.Errorf("the middle of the layer has alpha %d, want 255", a)
	}
}

// Every density, and nothing scaled up past what the source can fill.
func TestIconsAtEveryDensity(t *testing.T) {
	dir := t.TempDir()
	if err := WriteIcons(dir, source(t, dir, 512), canvas.Blue); err != nil {
		t.Fatal(err)
	}
	for _, d := range densities {
		for _, name := range []string{"ic_launcher.png", "ic_launcher_round.png"} {
			w, _, _ := read(t, filepath.Join(dir, "mipmap-"+d.name, name))
			if w != d.size {
				t.Errorf("%s/%s is %d across, want %d", d.name, name, w, d.size)
			}
		}
	}
	for _, name := range []string{"ic_launcher.xml", "ic_launcher_round.xml"} {
		if _, err := os.Stat(filepath.Join(dir, "mipmap-anydpi-v26", name)); err != nil {
			t.Errorf("no adaptive icon at %s", name)
		}
	}
	colours, err := os.ReadFile(filepath.Join(dir, "values", "colors.xml"))
	if err != nil {
		t.Fatal(err)
	}
	// canvas.Blue is 0xFF3E63DD.
	if want := "#FF3E63DD"; !contains(string(colours), want) {
		t.Errorf("the background colour is not %s:\n%s", want, colours)
	}
}

func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
