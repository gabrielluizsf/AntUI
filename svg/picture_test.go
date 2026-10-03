package svg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// The picture an <image> draws is read from where its href points — a data URI,
// an address over the network, or a file — fitted into the box it was given and
// painted like everything else, and what any of that cannot do is said out loud.

// solidPicture is a PNG of one colour, the picture most of these tests draw
// with: a colour that is easy to find again in a pixel and a size that fits
// the boxes below.
func solidPicture(t *testing.T, w, h int, c color.Color) []byte {
	t.Helper()
	return pictureAt(t, w, h, func(int, int) color.Color { return c })
}

// pictureAt is a PNG built pixel by pixel, for the pictures whose parts have
// to be told apart once they are on the canvas.
func pictureAt(t *testing.T, w, h int, at func(x, y int) color.Color) []byte {
	t.Helper()
	im := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			im.Set(x, y, at(x, y))
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, im); err != nil {
		t.Fatalf("encode the picture: %v", err)
	}
	return buf.Bytes()
}

// asDataURI is a picture carried in the href itself, the way a drawing on the
// web writes one small enough to fit in the file.
func asDataURI(picture []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(picture)
}

// mustParse reads a drawing and fails the test rather than carrying on without
// one, which is what every test here asks for first.
func mustParse(t *testing.T, src string) *Image {
	t.Helper()
	img, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return img
}

// wantColour is what a pixel came out as, named for what the test asked of it.
func wantColour(t *testing.T, cv *canvas.Canvas, x, y int, want canvas.Color) {
	t.Helper()
	if got := pixelAt(cv, x, y); got != want {
		t.Errorf("pixel (%d,%d) is %v, want %v", x, y, got, want)
	}
}

// wantNearColour is the same for a picture that went through a codec which
// does not promise every bit back — a JPEG, whose edges and blocks round a
// colour by a little however solid the block was.
func wantNearColour(t *testing.T, cv *canvas.Canvas, x, y int, r, g, b, a, tolerance int) {
	t.Helper()
	got := pixelAt(cv, x, y)
	for _, ch := range [][2]int{
		{int(got.R()), r}, {int(got.G()), g}, {int(got.B()), b}, {int(got.A()), a},
	} {
		if diff := ch[0] - ch[1]; diff > tolerance || diff < -tolerance {
			t.Errorf("pixel (%d,%d) is %v, want about (%d,%d,%d,%d) within %d", x, y, got, r, g, b, a, tolerance)
			return
		}
	}
}

func TestAnImagePaintsThePictureItsHrefPointsAt(t *testing.T) {
	// The whole trip: a picture carried in the address itself — with the line
	// break some encoders put down the middle of the base64, which is not part
	// of what it stands for — read into a canvas and laid over the box.
	pic := solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255})
	enc := base64.StdEncoding.EncodeToString(pic)
	mid := len(enc) / 2
	href := "data:image/png;base64," + enc[:mid] + "\n" + enc[mid:]
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4"/></svg>`, href))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a picture that was read said something was wrong: %s", got)
	}
	cv := img.Render(4, 4)
	wantColour(t, cv, 0, 0, canvas.RGB(255, 0, 0))
	wantColour(t, cv, 2, 2, canvas.RGB(255, 0, 0))
}

func TestAnImageReadsAJPEGAndAGIF(t *testing.T) {
	// The two other formats the standard library decodes, over the same box
	// and with the same tolerance for what a codec rounds.
	for _, tc := range []struct {
		name, media string
		encode      func(*image.RGBA) ([]byte, error)
	}{
		{"jpeg", "image/jpeg", func(im *image.RGBA) ([]byte, error) {
			var buf bytes.Buffer
			err := jpeg.Encode(&buf, im, nil)
			return buf.Bytes(), err
		}},
		{"gif", "image/gif", func(im *image.RGBA) ([]byte, error) {
			var buf bytes.Buffer
			// A GIF keeps the colours its palette holds, so the block is
			// written through a palette that holds the colour it went in as:
			// what comes back should be the colour it was, not the nearest the
			// default palette could guess.
			pal := image.NewPaletted(im.Bounds(), color.Palette{
				color.NRGBA{R: 255, G: 128, A: 255},
				color.NRGBA{},
			})
			draw.Draw(pal, pal.Bounds(), im, im.Bounds().Min, draw.Src)
			err := gif.Encode(&buf, pal, nil)
			return buf.Bytes(), err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			im := image.NewRGBA(image.Rect(0, 0, 8, 8))
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					im.Set(x, y, color.NRGBA{R: 255, G: 128, A: 255})
				}
			}
			data, err := tc.encode(im)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			href := "data:" + tc.media + ";base64," + base64.StdEncoding.EncodeToString(data)
			img := mustParse(t, fmt.Sprintf(
				`<svg viewBox="0 0 8 8"><image href="%s" x="0" y="0" width="8" height="8"/></svg>`, href))
			if got := img.Warnings().String(); got != "no warnings" {
				t.Errorf("a %s said something was wrong: %s", tc.name, got)
			}
			wantNearColour(t, img.Render(8, 8), 4, 4, 255, 128, 0, 255, 6)
		})
	}
}

func TestAnImageReadsAPictureCarriedAsPercentEscapes(t *testing.T) {
	// The other half of a data: URI: the bytes as text with what they stand
	// for written as %XX, which is how one is written by hand.
	pic := solidPicture(t, 4, 4, color.NRGBA{R: 255, B: 255, A: 255})
	var escaped strings.Builder
	for _, c := range pic {
		fmt.Fprintf(&escaped, "%%%02X", c)
	}
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="data:image/png,%s" x="0" y="0" width="4" height="4"/></svg>`,
		escaped.String()))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a percent-escaped picture said something was wrong: %s", got)
	}
	wantColour(t, img.Render(4, 4), 2, 2, canvas.RGB(255, 0, 255))
}

func TestAnImageReadsADataURIWithItsPaddingLeftOff(t *testing.T) {
	// Some encoders leave the `=` padding off the end of the base64, which is
	// not what the address is written with but is what arrives soon enough to
	// be worth reading rather than refusing.
	pic := solidPicture(t, 4, 4, color.NRGBA{G: 255, A: 255})
	href := "data:image/png;base64," + base64.RawStdEncoding.EncodeToString(pic)
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4"/></svg>`, href))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("unpadded base64 said something was wrong: %s", got)
	}
	wantColour(t, img.Render(4, 4), 2, 2, canvas.RGB(0, 255, 0))
}

func TestAnImageReadsAPictureFromAFile(t *testing.T) {
	// A picture beside the drawing, which is how one is written while it is
	// being made: the address is the file it names.
	pic := solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255})
	path := filepath.Join(t.TempDir(), "picture.png")
	if err := os.WriteFile(path, pic, 0o644); err != nil {
		t.Fatal(err)
	}
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4"/></svg>`, path))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a picture on disk said something was wrong: %s", got)
	}
	wantColour(t, img.Render(4, 4), 2, 2, canvas.RGB(255, 0, 0))
}

func TestAnImageFetchesAPictureOverHTTPAndOnlyOnce(t *testing.T) {
	// A CDN address is fetched while the drawing is read, and the same address
	// twice in one drawing is fetched once: the picture belongs to the address,
	// not to the element that named it.
	pic := solidPicture(t, 4, 4, color.NRGBA{G: 255, A: 255})
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "image/png")
		w.Write(pic)
	}))
	defer srv.Close()
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 8 4">
		<image href="%[1]s/pic.png" x="0" y="0" width="4" height="4"/>
		<image href="%[1]s/pic.png" x="4" y="0" width="4" height="4"/>
	</svg>`, srv.URL))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a picture over http said something was wrong: %s", got)
	}
	cv := img.Render(8, 4)
	wantColour(t, cv, 1, 1, canvas.RGB(0, 255, 0))
	wantColour(t, cv, 5, 1, canvas.RGB(0, 255, 0))
	if n := hits.Load(); n != 1 {
		t.Errorf("the same picture was fetched %d times, want 1", n)
	}
}

func TestAPictureTheServerRefusesIsSaidOutLoud(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s/pic.png" width="4" height="4"/></svg>`, srv.URL))
	if got := img.Warnings().String(); !strings.Contains(got, "404") {
		t.Errorf("a refused picture was not said out loud: %s", got)
	}
}

func TestProtocolRelativeAddressesAreReadAsHTTPS(t *testing.T) {
	// There is no page above a drawing to take a scheme from, so a
	// protocol-relative address is the https it means wherever it is opened.
	for _, tc := range []struct {
		href, want string
		ok         bool
	}{
		{"//cdn.example/pic.png", "https://cdn.example/pic.png", true},
		{"http://example/pic.png", "http://example/pic.png", true},
		{"https://example/pic.png", "https://example/pic.png", true},
		{"pic.png", "", false},
		{"data:image/png;base64,AAAA", "", false},
	} {
		got, ok := pictureHTTPURL(tc.href)
		if ok != tc.ok || got != tc.want {
			t.Errorf("pictureHTTPURL(%q) = %q, %v; want %q, %v", tc.href, got, ok, tc.want, tc.ok)
		}
	}
}

func TestAnImageFallsBackToTheXlinkHref(t *testing.T) {
	// The spelling a drawing written for SVG 1.1 carries, which still turns up
	// beside every modern file that has one.
	pic := solidPicture(t, 4, 4, color.NRGBA{R: 255, G: 255, A: 255})
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image xlink:href="%s" x="0" y="0" width="4" height="4"/></svg>`,
		asDataURI(pic)))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("an xlink:href said something was wrong: %s", got)
	}
	wantColour(t, img.Render(4, 4), 2, 2, canvas.RGB(255, 255, 0))
}

func TestAnImageTakesItsBoxFromTheDrawing(t *testing.T) {
	// Lengths in the drawing's own units and percentages of the drawing, side
	// by side: `x` and `y` absent are the origin, and a percentage is of the
	// viewBox the same one a gradient's is.
	red := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	green := asDataURI(solidPicture(t, 4, 4, color.NRGBA{G: 255, A: 255}))
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 100 100">
		<image href="%[1]s" x="10" y="10" width="20" height="20"/>
		<image href="%[2]s" x="50%%" y="50%%" width="25%%" height="25%%"/>
	</svg>`, red, green))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("boxes in units and percentages said something was wrong: %s", got)
	}
	cv := img.Render(100, 100)
	wantColour(t, cv, 15, 15, canvas.RGB(255, 0, 0))
	wantColour(t, cv, 60, 60, canvas.RGB(0, 255, 0))
	wantColour(t, cv, 40, 40, canvas.Transparent)
	wantColour(t, cv, 90, 90, canvas.Transparent)
}

func TestAnImageMeetsItsPictureInTheBox(t *testing.T) {
	// A wide picture in a square box keeps its own shape and is centred in it
	// by default; the alignment moves it to the corner that was asked for and
	// what is left over of the box stays empty.
	pic := solidPicture(t, 4, 2, color.NRGBA{B: 255, A: 255})
	for _, tc := range []struct {
		preserve    string
		full, empty int
	}{
		{"", 1, 3},
		{"xMinYMin meet", 0, 2},
		{"xMaxYMax meet", 2, 0},
	} {
		img := mustParse(t, fmt.Sprintf(
			`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4" preserveAspectRatio="%s"/></svg>`,
			asDataURI(pic), tc.preserve))
		if got := img.Warnings().String(); got != "no warnings" {
			t.Errorf("%q said something was wrong: %s", tc.preserve, got)
		}
		cv := img.Render(4, 4)
		wantColour(t, cv, 2, tc.full, canvas.RGB(0, 0, 255))
		wantColour(t, cv, 2, tc.empty, canvas.Transparent)
	}
}

func TestADeferredPreserveAspectRatioIsReadWithoutIt(t *testing.T) {
	// SVG 1.1 writes a leading `defer` for a picture pointed at from a
	// drawing with a fit of its own; there is no drawing behind an address
	// here, so the word is passed over and the rest is read as it stands.
	pic := solidPicture(t, 4, 2, color.NRGBA{B: 255, A: 255})
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4" preserveAspectRatio="defer xMinYMin meet"/></svg>`,
		asDataURI(pic)))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a deferred preserveAspectRatio said something was wrong: %s", got)
	}
	cv := img.Render(4, 4)
	wantColour(t, cv, 2, 0, canvas.RGB(0, 0, 255))
	wantColour(t, cv, 2, 3, canvas.Transparent)
}

func TestAnImageSlicesItsPictureToTheBox(t *testing.T) {
	// A picture cut down to the box: the left half red and the right half
	// blue, so which half the alignment pins is there to see — and the whole
	// box covered however, which a meet would have left clear at the top and
	// the bottom.
	pic := pictureAt(t, 4, 2, func(x, y int) color.Color {
		if x < 2 {
			return color.NRGBA{R: 255, A: 255}
		}
		return color.NRGBA{B: 255, A: 255}
	})
	for _, tc := range []struct {
		preserve    string
		left, right canvas.Color
	}{
		{"xMinYMin slice", canvas.RGB(255, 0, 0), canvas.RGB(255, 0, 0)},
		{"xMaxYMax slice", canvas.RGB(0, 0, 255), canvas.RGB(0, 0, 255)},
		{"xMidYMid slice", canvas.RGB(255, 0, 0), canvas.RGB(0, 0, 255)},
	} {
		img := mustParse(t, fmt.Sprintf(
			`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4" preserveAspectRatio="%s"/></svg>`,
			asDataURI(pic), tc.preserve))
		if got := img.Warnings().String(); got != "no warnings" {
			t.Errorf("%q said something was wrong: %s", tc.preserve, got)
		}
		cv := img.Render(4, 4)
		wantColour(t, cv, 1, 0, tc.left)
		wantColour(t, cv, 1, 3, tc.left)
		wantColour(t, cv, 3, 3, tc.right)
	}
}

func TestAnImageStretchesItsPictureWithNone(t *testing.T) {
	// The fit that throws the picture's own shape away for the box: a wide
	// picture over every pixel of a square one, top row and bottom included.
	pic := solidPicture(t, 4, 2, color.NRGBA{G: 255, A: 255})
	img := mustParse(t, fmt.Sprintf(
		`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4" preserveAspectRatio="none"/></svg>`,
		asDataURI(pic)))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a stretched picture said something was wrong: %s", got)
	}
	cv := img.Render(4, 4)
	wantColour(t, cv, 2, 0, canvas.RGB(0, 255, 0))
	wantColour(t, cv, 2, 3, canvas.RGB(0, 255, 0))
	wantColour(t, cv, 0, 2, canvas.RGB(0, 255, 0))
}

func TestAnImageDrawsTheOpacityItWasGivenWithoutFadingTheRest(t *testing.T) {
	// The opacity goes into pixels of this element's own, so the second one
	// pointing at the same address comes out whole: what the two share is the
	// picture as it was written, not what either of them does with it.
	pic := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 8 4">
		<image href="%[1]s" x="0" y="0" width="4" height="4" opacity="0.5"/>
		<image href="%[1]s" x="4" y="0" width="4" height="4"/>
	</svg>`, pic))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("two pictures at two opacities said something was wrong: %s", got)
	}
	cv := img.Render(8, 4)
	faint := pixelAt(cv, 2, 2)
	if a := int(faint.A()); a < 126 || a > 130 {
		t.Errorf("the faint picture has alpha %d, want about 128", a)
	}
	if r := int(faint.R()); r < 250 {
		t.Errorf("the faint picture came out with red %d, want it still red", r)
	}
	wantColour(t, cv, 6, 2, canvas.RGB(255, 0, 0))
}

func TestAnImageKeepsItsColourThroughTheAlphaTheDecoderReturns(t *testing.T) {
	// The decoder hands back channels already multiplied by their alpha, and a
	// canvas holds them straight: half a red pixel read back as it stands would
	// come out a dark red wherever it was laid over.
	for _, tc := range []struct {
		name  string
		build func() image.Image
	}{
		{"8-bit", func() image.Image {
			im := image.NewNRGBA(image.Rect(0, 0, 4, 4))
			for y := 0; y < 4; y++ {
				for x := 0; x < 4; x++ {
					im.Set(x, y, color.NRGBA{R: 255, A: 128})
				}
			}
			return im
		}},
		{"16-bit", func() image.Image {
			im := image.NewNRGBA64(image.Rect(0, 0, 4, 4))
			for y := 0; y < 4; y++ {
				for x := 0; x < 4; x++ {
					im.Set(x, y, color.NRGBA64{R: 65535, A: 32768})
				}
			}
			return im
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := png.Encode(&buf, tc.build()); err != nil {
				t.Fatalf("encode the picture: %v", err)
			}
			img := mustParse(t, fmt.Sprintf(
				`<svg viewBox="0 0 4 4"><image href="%s" x="0" y="0" width="4" height="4"/></svg>`,
				asDataURI(buf.Bytes())))
			if got := img.Warnings().String(); got != "no warnings" {
				t.Errorf("a translucent picture said something was wrong: %s", got)
			}
			wantNearColour(t, img.Render(4, 4), 2, 2, 255, 0, 0, 128, 2)
		})
	}
}

func TestAPictureThatIsNotPartOfThePictureIsNeverFetched(t *testing.T) {
	// An address nothing is going to show is not worth reading: an element
	// turned off, one still waiting in a definition, and one asked for at an
	// opacity of nothing all leave the file alone — and say nothing about an
	// address that might well be broken, since nothing was asked of it.
	for name, src := range map[string]string{
		"display none": `<svg viewBox="0 0 4 4"><image style="display:none" href="missing.png" width="4" height="4"/></svg>`,
		"in a defs":    `<svg viewBox="0 0 4 4"><defs><image href="missing.png" width="4" height="4"/></defs></svg>`,
		"opacity zero": `<svg viewBox="0 0 4 4"><image href="missing.png" width="4" height="4" opacity="0"/></svg>`,
	} {
		t.Run(name, func(t *testing.T) {
			img := mustParse(t, src)
			if n := len(img.Warnings()); n != 0 {
				t.Errorf("it said %d things about a picture nothing asked for: %s", n, img.Warnings())
			}
		})
	}
}

func TestAPictureItCannotReadIsSaidOutLoud(t *testing.T) {
	good := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	for _, tc := range []struct{ name, src, want string }{
		{"no href",
			`<svg viewBox="0 0 4 4"><image width="4" height="4"/></svg>`,
			"no href"},
		{"no room",
			`<svg viewBox="0 0 4 4"><image href="pic.png"/></svg>`,
			"no width or height"},
		{"a width of nothing",
			`<svg viewBox="0 0 4 4"><image href="pic.png" width="0" height="4"/></svg>`,
			"no width or height"},
		{"a length that is not one",
			`<svg viewBox="0 0 4 4"><image href="pic.png" width="wide" height="4"/></svg>`,
			"not a length"},
		{"a percentage of nothing",
			`<svg viewBox="0 0 4 4"><image href="pic.png" width="%" height="4"/></svg>`,
			"not a length"},
		{"an escape that stands for nothing",
			`<svg viewBox="0 0 4 4"><image href="data:image/png,%ZZ" width="4" height="4"/></svg>`,
			"could not be read"},
		{"a file that is not there",
			`<svg viewBox="0 0 4 4"><image href="pic.png" width="4" height="4"/></svg>`,
			"could not be read"},
		{"bytes that are not a picture",
			`<svg viewBox="0 0 4 4"><image href="data:text/plain,hello" width="4" height="4"/></svg>`,
			"could not be read"},
		{"an address with no comma",
			`<svg viewBox="0 0 4 4"><image href="data:image/png;base64" width="4" height="4"/></svg>`,
			"could not be read"},
		{"base64 that is not",
			`<svg viewBox="0 0 4 4"><image href="data:image/png;base64,!!!!" width="4" height="4"/></svg>`,
			"could not be read"},
		{"an alignment that does not exist", fmt.Sprintf(
			`<svg viewBox="0 0 4 4"><image href="%s" width="4" height="4" preserveAspectRatio="sideways meet"/></svg>`, good),
			"is not an alignment"},
		{"a meet or a slice that is not", fmt.Sprintf(
			`<svg viewBox="0 0 4 4"><image href="%s" width="4" height="4" preserveAspectRatio="xMidYMid round"/></svg>`, good),
			"where a meet or a slice goes"},
		{"more than it holds", fmt.Sprintf(
			`<svg viewBox="0 0 4 4"><image href="%s" width="4" height="4" preserveAspectRatio="xMidYMid meet tail"/></svg>`, good),
			"the rest is left out"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			img := mustParse(t, tc.src)
			if got := img.Warnings().String(); !strings.Contains(got, tc.want) {
				t.Errorf("the warnings say:\n%s\nand should contain %q", got, tc.want)
			}
		})
	}
}

func TestAPictureThroughAUseGoesWhereTheUsePutsIt(t *testing.T) {
	// A definition that is a photograph rather than a shape: built where the
	// `<use>` points at it, moved by the x and y the `<use>` was given, the
	// same as any other definition this drawing keeps.
	pic := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 8 4">
		<defs><image id="p" href="%s" x="0" y="0" width="4" height="4"/></defs>
		<use href="#p" x="4"/>
	</svg>`, pic))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a picture through a use said something was wrong: %s", got)
	}
	cv := img.Render(8, 4)
	wantColour(t, cv, 1, 1, canvas.Transparent)
	wantColour(t, cv, 5, 1, canvas.RGB(255, 0, 0))
}

func TestAPictureIsCutByAClipLikeAnythingElse(t *testing.T) {
	// The clip goes over the picture the same way it goes over a shape: the
	// picture is drawn whole into its element's picture first and the picture
	// is cut once, the left half of the box kept here.
	pic := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 4 4">
		<clipPath id="c"><rect x="0" y="0" width="2" height="4"/></clipPath>
		<image href="%s" x="0" y="0" width="4" height="4" clip-path="url(#c)"/>
	</svg>`, pic))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a clipped picture said something was wrong: %s", got)
	}
	cv := img.Render(4, 4)
	wantColour(t, cv, 1, 1, canvas.RGB(255, 0, 0))
	wantColour(t, cv, 3, 1, canvas.Transparent)
}

func TestAMaskMeasuredAgainstTheBoxOfAPictureCutsIt(t *testing.T) {
	// A mask whose region is written against the box the picture covers: the
	// box is what the picture has to show — it has no outline to measure the
	// way a shape does — and the half of it the region keeps is the half that
	// comes out.
	pic := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 4 4">
		<mask id="m" maskUnits="objectBoundingBox" x="0" y="0" width="0.5" height="1">
			<rect x="0" y="0" width="4" height="4" fill="white"/>
		</mask>
		<image href="%s" x="0" y="0" width="4" height="4" mask="url(#m)"/>
	</svg>`, pic))
	if got := img.Warnings().String(); got != "no warnings" {
		t.Errorf("a masked picture said something was wrong: %s", got)
	}
	cv := img.Render(4, 4)
	wantColour(t, cv, 1, 1, canvas.RGB(255, 0, 0))
	wantColour(t, cv, 3, 1, canvas.Transparent)
}

func TestAPictureInsideAClipPathSaysItHasNoOutline(t *testing.T) {
	// A clip is cut with outlines, and a photograph has none: the same thing
	// the writing inside a clipPath says, beside the same nothing it leaves
	// behind in the clip.
	pic := asDataURI(solidPicture(t, 4, 4, color.NRGBA{R: 255, A: 255}))
	img := mustParse(t, fmt.Sprintf(`<svg viewBox="0 0 4 4">
		<clipPath id="c"><image href="%s" width="4" height="4"/></clipPath>
		<rect x="0" y="0" width="4" height="4" clip-path="url(#c)"/>
	</svg>`, pic))
	got := img.Warnings().String()
	if !strings.Contains(got, "no outline to cut with") {
		t.Errorf("a picture in a clipPath was not said out loud: %s", got)
	}
	if !strings.Contains(got, "image") {
		t.Errorf("the warning does not name the element it came from: %s", got)
	}
}
