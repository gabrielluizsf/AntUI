package apk

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui/canvas"
)

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
	if err := writeDensities(dir, src); err != nil {
		return err
	}
	if err := writeAdaptive(dir); err != nil {
		return err
	}
	return writeColors(dir, background)
}
