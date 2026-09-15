package apk

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
)

// The launcher icon at each density, in pixels. They are 48 device-
// independent pixels multiplied out, which is the size every launcher asks
// for and has since the beginning.
var densities = []struct {
	name string
	size int
	// adaptive is the same 48dp icon's *adaptive* layer, which is 108dp —
	// larger than what is shown, because the launcher crops it to whatever
	// shape that device uses and may move it about while it does.
	adaptive int
}{
	{"mdpi", 48, 108},
	{"hdpi", 72, 162},
	{"xhdpi", 96, 216},
	{"xxhdpi", 144, 324},
	{"xxxhdpi", 192, 432},
}

// safeFraction is how much of an adaptive layer is certain to be visible.
// The layer is 108dp and the guaranteed circle inside it is 72dp, so
// anything that must be seen has to be within two thirds of the middle —
// the rest may be cropped, and on some launchers is.
const safeFraction = 72.0 / 108.0

// WriteIcons turns one picture into every file the launcher wants, under
// dir: a bitmap at each density, a round one, and the adaptive pair that
// every device since Android 8 draws.
//
// It is exported because it is useful on its own — "antuiapk icon" is this
// function and nothing else — and because an app that has its own build
// wants the same set of files without the rest of the package.
//
// It uses antui's own scaler, which is an alpha-weighted box filter: every
// pixel of the result is the average of the block of the original that lands
// on it, weighted so that the colour of a transparent corner does not bleed
// into the edge of the shape. An icon is small and square and comes from
// something large and square, which is the one case that filter is exactly
// right for.
func WriteIcons(dir, source string, background canvas.Color) error {
	src, err := loadCanvas(source)
	if err != nil {
		return err
	}
	if src.Width < 192 || src.Height < 192 {
		// Not an error — a small source still makes every size — but a
		// launcher icon at 192 is what a modern phone shows, and scaling up
		// to it gains nothing.
		fmt.Fprintf(os.Stderr,
			"antuiapk: the icon is %dx%d; 512x512 is what the store asks for\n",
			src.Width, src.Height)
	}

	for _, d := range densities {
		folder := filepath.Join(dir, "mipmap-"+d.name)
		if err := os.MkdirAll(folder, 0o755); err != nil {
			return err
		}
		square := canvas.IconScaled(src, d.size)
		if square == nil {
			return fmt.Errorf("apk: cannot scale the icon to %d", d.size)
		}
		if err := writePNG(filepath.Join(folder, "ic_launcher.png"), square); err != nil {
			return err
		}
		// The round one is the same picture with the corners taken off. A
		// launcher that asks for it is one that will not mask what it is
		// given, so the masking has to be done here.
		if err := writePNG(filepath.Join(folder, "ic_launcher_round.png"),
			rounded(square)); err != nil {
			return err
		}
		if err := writePNG(filepath.Join(folder, "ic_launcher_foreground.png"),
			adaptiveLayer(src, d.adaptive)); err != nil {
			return err
		}
	}

	// The adaptive icon itself, which is a description rather than a
	// picture: a background, a foreground, and — from Android 13 — a
	// monochrome layer the system tints to match the user's wallpaper.
	// Anydpi-v26 means "every density, from API 26", so it wins over the
	// PNGs above wherever adaptive icons exist and is ignored where they do
	// not.
	folder := filepath.Join(dir, "mipmap-anydpi-v26")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return err
	}
	adaptive := `<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background" />
    <foreground android:drawable="@mipmap/ic_launcher_foreground" />
    <monochrome android:drawable="@mipmap/ic_launcher_foreground" />
</adaptive-icon>
`
	for _, name := range []string{"ic_launcher.xml", "ic_launcher_round.xml"} {
		if err := os.WriteFile(filepath.Join(folder, name), []byte(adaptive), 0o644); err != nil {
			return err
		}
	}

	values := filepath.Join(dir, "values")
	if err := os.MkdirAll(values, 0o755); err != nil {
		return err
	}
	colors := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <color name="ic_launcher_background">#%02X%02X%02X%02X</color>
</resources>
`, background.A(), background.R(), background.G(), background.B())
	return os.WriteFile(filepath.Join(values, "colors.xml"), []byte(colors), 0o644)
}

// clear empties a canvas to nothing at all.
//
// Not Canvas.Clear, which forces the alpha to 255 — deliberately, because it
// is for clearing a *window* and a frame on screen is opaque. An icon is the
// opposite: what is not drawn has to be see-through, and a rounded corner
// filled with opaque black is a black square.
func clear(cv *canvas.Canvas) {
	for y := range cv.Height {
		for x := range cv.Width {
			cv.Put(x, y, canvas.Transparent)
		}
	}
}

// adaptiveLayer makes the foreground of an adaptive icon: the picture at two
// thirds of the layer, centred, with nothing around it.
func adaptiveLayer(src *canvas.Canvas, size int) *canvas.Canvas {
	out, err := canvas.NewCanvas(size, size)
	if err != nil {
		return nil
	}
	clear(out)
	inner := int(float64(size) * safeFraction)
	scaled := canvas.IconScaled(src, inner)
	if scaled == nil {
		return out
	}
	at := (size - inner) / 2
	out.Blit(at, at, scaled)
	return out
}

// rounded takes the corners off a square icon.
func rounded(src *canvas.Canvas) *canvas.Canvas {
	out, err := canvas.NewCanvas(src.Width, src.Height)
	if err != nil {
		return src
	}
	clear(out)
	// The radius is half the width, and the comparison is in squared
	// distance so that there is no square root per pixel.
	radius := src.Width / 2
	limit := radius * radius
	for y := range src.Height {
		dy := y - radius
		for x := range src.Width {
			dx := x - radius
			if dx*dx+dy*dy <= limit {
				out.Put(x, y, src.At(x, y))
			}
		}
	}
	return out
}

// loadCanvas reads a PNG or a JPEG into a antui canvas.
func loadCanvas(path string) (*canvas.Canvas, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("apk: cannot read the icon: %w", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("apk: %s is not a PNG or a JPEG: %w", path, err)
	}
	b := img.Bounds()
	cv, err := canvas.NewCanvas(b.Dx(), b.Dy())
	if err != nil {
		return nil, err
	}
	for y := range b.Dy() {
		for x := range b.Dx() {
			r, g, bl, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			// RGBA gives alpha-premultiplied 16-bit channels; a Canvas holds
			// straight 8-bit ones, so the multiplication has to be undone or
			// everything transparent comes out too dark.
			if a == 0 {
				cv.Put(x, y, canvas.Transparent)
				continue
			}
			cv.Put(x, y, canvas.RGBA(
				uint8(r*0xFF/a), uint8(g*0xFF/a), uint8(bl*0xFF/a), uint8(a>>8)))
		}
	}
	return cv, nil
}

// writePNG saves a canvas.
func writePNG(path string, cv *canvas.Canvas) error {
	if cv == nil {
		return fmt.Errorf("apk: nothing to write to %s", path)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, canvasImage{cv})
}

// canvasImage is a canvas.Canvas seen as an image.Image, which is what the
// standard library's encoder wants. Canvas has At, but it answers a
// canvas.Color rather than a color.Color, so it cannot be the interface
// itself.
type canvasImage struct{ cv *canvas.Canvas }

func (c canvasImage) ColorModel() color.Model { return color.NRGBAModel }
func (c canvasImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.cv.Width, c.cv.Height)
}
func (c canvasImage) At(x, y int) color.Color {
	p := c.cv.At(x, y)
	return color.NRGBA{R: p.R(), G: p.G(), B: p.B(), A: p.A()}
}
