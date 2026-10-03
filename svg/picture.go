package svg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"  // a GIF is one of the three pictures this package reads
	_ "image/jpeg" // a JPEG is one of the three pictures this package reads
	_ "image/png"  // a PNG is one of the three pictures this package reads
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// Picture is the bitmap an `<image>` draws and the box it is drawn into. It is
// the whole of what was asked for: where the picture came from has already been
// read by the time one of these exists, so painting needs nothing but the
// pixels and the four numbers that say where they land.
type Picture struct {
	// Pixels is the picture itself, straight-alpha the way a canvas holds it,
	// at the size it was encoded at. Two `<image>`s pointing at the same
	// address share these pixels and are never written to: fitting is done with
	// where they are put rather than by changing them.
	Pixels *canvas.Canvas
	// Box is the rectangle of the drawing the picture is fitted into, as x, y,
	// width and height in the drawing's own coordinates, with a percentage
	// already counted as its fraction of the drawing. It is as the file wrote
	// it, before any transform on the element moves it — see [paintPicture],
	// which puts the transform on afterwards the same way writing takes it.
	Box [4]float64
	// Fit is whether the picture keeps its own shape inside the box (meet), is
	// cut down to the box (slice), or is stretched to fill it (none), which is
	// what `preserveAspectRatio` says.
	Fit PictureFit
	// Align is which corner or edge of the box the picture is pinned to when
	// the fit leaves anything over: one of the nine alignments the attribute
	// names, `xMidYMid` being the one a drawing that says nothing gets.
	Align PictureAlign
}

// PictureFit is what `preserveAspectRatio` says about the picture's own shape:
// kept whole inside the box, cut down to the box, or thrown away for the box.
type PictureFit uint8

const (
	// PictureFitMeet keeps the whole picture inside the box, padding what is
	// left over with nothing. It is the default and the usual way a photograph
	// is shown.
	PictureFitMeet PictureFit = iota
	// PictureFitSlice covers the box with the picture, cutting off what falls
	// outside it — the way a full-bleed background is shown.
	PictureFitSlice
	// PictureFitNone stretches the picture to the box with no regard for the
	// shape it was made in.
	PictureFitNone
)

// PictureAlign is which of the nine points of the box the picture is pinned to
// when the fit leaves something over: an edge or the middle of each axis,
// spelled in the file as `xMinYMid` and its eight siblings. It is read as the
// one number it maps to — the axis positions come out of it with [PictureAlign.offsets].
type PictureAlign uint8

const (
	AlignXMinYMin PictureAlign = iota
	AlignXMidYMin
	AlignXMaxYMin
	AlignXMinYMid
	AlignXMidYMid
	AlignXMaxYMid
	AlignXMinYMax
	AlignXMidYMax
	AlignXMaxYMax
)

// offsets is how far along each side of the box the picture starts: 0 for the
// start edge, a half for the middle, a whole one for the far edge — the three
// positions each axis of the nine names has.
func (a PictureAlign) offsets() (fx, fy float64) {
	return float64(uint(a)%3) / 2, float64(uint(a)/3) / 2
}

// pictureAligns is every alignment the attribute spells, lowercased because the
// values are matched case-insensitively the same way the names are.
var pictureAligns = map[string]PictureAlign{
	"xminymin": AlignXMinYMin, "xmidymin": AlignXMidYMin, "xmaxymin": AlignXMaxYMin,
	"xminymid": AlignXMinYMid, "xmidymid": AlignXMidYMid, "xmaxymid": AlignXMaxYMid,
	"xminymax": AlignXMinYMax, "xmidymax": AlignXMidYMax, "xmaxymax": AlignXMaxYMax,
}

// imageNode is the `<image>` read whole: where the picture comes from, the box
// it goes into, how it is fitted there, and the opacity it is drawn at. It
// answers the node to paint, with nothing on it where there was nothing to
// paint — a drawing that asked for a picture this cannot fetch or read leaves
// the hole and says which it was, the same way a `url(#id)` that is not there
// leaves a shape unfilled.
//
// The picture is fetched while the drawing is read rather than when it is
// painted, because reading is where every other warning this package gives
// comes from and where a drawing is shared: what went wrong is said once, to
// whoever called Parse, instead of once a frame for as long as the drawing is
// on screen.
func (img *Image) imageNode(e *element, st Style, warn func(string, ...any)) *Node {
	n := &Node{Name: e.Name, Style: st}
	href := e.attr("href")
	if href == "" {
		// The older spelling of the same attribute, which is the one a drawing
		// written for SVG 1.1 carries — the same fallback a `<use>` takes.
		href = e.attr("xlink:href")
	}
	if href == "" {
		warn("it has no href, so there is no picture for it to draw")
		return n
	}
	box := [4]float64{
		img.pictureLength(warn, e, "x", 0),
		img.pictureLength(warn, e, "y", 1),
		img.pictureLength(warn, e, "width", 0),
		img.pictureLength(warn, e, "height", 1),
	}
	if box[2] <= 0 || box[3] <= 0 {
		// A box of no area is where the picture would go and there is no there,
		// the same as a rectangle with no width: nothing is drawn, and this
		// says so rather than the drawing quietly keeping a picture nobody can
		// see. Fetching first would read a file only to throw it away.
		warn("it has no width or height to draw the picture into, so the picture is left out")
		return n
	}
	fit, align := readPreserveAspect(e.attr("preserveAspectRatio"), warn)
	if st.Opacity <= 0 {
		// Opacity of nothing: the picture would come out transparent however
		// good it is, so it is not worth fetching. Unlike display:none this is
		// not about hiding the element from the drawing — it is already not
		// going to be seen, so there is nothing to say about it either.
		return n
	}
	pixels := img.pictureFor(href, warn)
	if pixels == nil {
		return n
	}
	if st.Opacity < 1 {
		// The opacity is put into pixels of their own rather than onto the
		// ones in the cache, which every other element pointing at the same
		// address shares: what is in the cache is the picture as it was
		// written, and one element drawing it faint does not fade the rest.
		pixels = fadePicture(pixels, st.Opacity)
	}
	n.Pic = &Picture{Pixels: pixels, Box: box, Fit: fit, Align: align}
	return n
}

// pictureFor is the picture an href points at, read once and kept under that
// address for every other `<image>` in the same drawing. A picture that cannot
// be read says so and answers nothing, which is what leaves the hole in the
// picture with a warning beside it.
func (img *Image) pictureFor(href string, warn func(string, ...any)) *canvas.Canvas {
	if pic := img.pictures[href]; pic != nil {
		return pic
	}
	data, err := fetchPicture(href)
	if err == nil {
		var pic *canvas.Canvas
		pic, err = decodePicture(data)
		if err == nil {
			if img.pictures == nil {
				img.pictures = map[string]*canvas.Canvas{}
			}
			img.pictures[href] = pic
			return pic
		}
	}
	warn("the picture %q could not be read: %v", href, err)
	return nil
}

// pictureLength is one of the four numbers the box is written with: a length
// in the drawing's own units, or a percentage of the drawing — which is the
// viewBox, the same thing a percentage of a gradient is a fraction of. An
// attribute that is not there is zero, and one that is written as something
// this cannot read says so and counts as nothing, the same as every other
// length the drawing writes. axis is which side of the drawing a percentage is
// a fraction of: 0 for across, 1 for down.
func (img *Image) pictureLength(warn func(string, ...any), e *element, attr string, axis int) float64 {
	if !e.hasAttr(attr) {
		return 0
	}
	raw := strings.TrimSpace(e.attr(attr))
	if p, ok := strings.CutSuffix(raw, "%"); ok {
		v, err := parseNumber(p)
		if err != nil {
			warn("the %s %q is not a length, so it is taken as none", attr, e.attr(attr))
			return 0
		}
		span := img.ViewBox[2]
		if axis == 1 {
			span = img.ViewBox[3]
		}
		return v / 100 * span
	}
	v, ok := parseLength(raw)
	if !ok {
		warn("the %s %q is not a length, so it is taken as none", attr, e.attr(attr))
	}
	return v
}

// readPreserveAspect is what the attribute says about the picture's own shape
// inside its box: an alignment, then whether the picture is met or sliced —
// either of which may be left out, falling back to the centred meeting a
// drawing that says nothing gets. Anything this cannot read says so and leaves
// the default standing, because a picture fitted the wrong way is one an author
// can see, while a picture fitted the default way beside a warning is one they
// can go and fix.
func readPreserveAspect(raw string, warn func(string, ...any)) (PictureFit, PictureAlign) {
	fit, align := PictureFitMeet, AlignXMidYMid
	fields := strings.Fields(raw)
	// SVG 1.1 allows a leading `defer`, which asks for the fit of a drawing
	// this picture is pointed at rather than of this one. There is no drawing
	// behind an address being fetched here, so the word is passed over and the
	// rest is read as it would be without it.
	i := 0
	if i < len(fields) && strings.EqualFold(fields[i], "defer") {
		i++
	}
	if i < len(fields) {
		switch tok := strings.ToLower(fields[i]); {
		case tok == "none":
			fit = PictureFitNone
		default:
			if a, ok := pictureAligns[tok]; ok {
				align = a
			} else {
				warn("the preserveAspectRatio %q begins with %q, which is not an alignment, so the picture is met and centred", raw, fields[i])
			}
		}
		i++
	}
	if fit != PictureFitNone && i < len(fields) {
		// With no alignment there is nothing for a meet or a slice to keep the
		// shape of — a stretched picture meets and slices the same way — so the
		// second word only counts where the first one left a shape to keep.
		switch strings.ToLower(fields[i]) {
		case "meet":
		case "slice":
			fit = PictureFitSlice
		default:
			warn("the preserveAspectRatio %q has %q where a meet or a slice goes, so the picture is met", raw, fields[i])
		}
		i++
	}
	if i < len(fields) {
		warn("the preserveAspectRatio %q has more in it than an alignment and a meet or a slice, so the rest is left out", raw)
	}
	return fit, align
}

// maxPictureBytes is how big a picture fetched over the network may be before
// it is refused rather than read: a server that answers with something the
// size of a disk should not be able to fill memory before the decode ever
// gets a say, and a truncated picture fails to decode with a message of its
// own anyway, so the refusal at the door is the clearer of the two.
const maxPictureBytes = 64 << 20

// pictureHTTP is what fetches a picture from an address, with a timeout: a
// library that reads a drawing synchronously cannot wait for a server for
// ever, and a picture that does not arrive in ten seconds is a picture that is
// not going to arrive in a way this drawing can use.
var pictureHTTP = &http.Client{Timeout: 10 * time.Second}

// fetchPicture is the bytes behind an href: a `data:` URI unwrapped where it
// stands, an address fetched over the network — a protocol-relative one being
// the https it is meant to be, since there is no page here to take a scheme
// from — and anything else read as the file it names, which is what makes a
// picture beside the drawing work while it is being written.
func fetchPicture(href string) ([]byte, error) {
	if strings.HasPrefix(href, "data:") {
		return dataURI(href)
	}
	if target, ok := pictureHTTPURL(href); ok {
		return getPicture(target)
	}
	return os.ReadFile(href)
}

// pictureHTTPURL is the address a href is fetched from where it is one at all:
// http and https as they stand, and a protocol-relative `//host/path` taken to
// https, which is the scheme a drawing with no page above it means. Anything
// else — a data URI, a file path — is not an address this function names.
func pictureHTTPURL(href string) (string, bool) {
	switch {
	case strings.HasPrefix(href, "http://"), strings.HasPrefix(href, "https://"):
		return href, true
	case strings.HasPrefix(href, "//"):
		return "https:" + href, true
	}
	return "", false
}

// getPicture is one address read whole, or the reason it could not be: a
// server that did not answer with the picture is said by the answer it did
// give, and a body past the size this will take is refused before it is in
// memory rather than after.
func getPicture(target string) ([]byte, error) {
	resp, err := pictureHTTP.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("it answered %s", resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxPictureBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxPictureBytes {
		return nil, fmt.Errorf("it is bigger than %d bytes", maxPictureBytes)
	}
	return data, nil
}

// dataURI is the bytes a `data:` address carries: the media type it announces
// is left for the decoder to have an opinion about, base64 is unwrapped when
// the address says so — with or without the padding some encoders leave off —
// and anything else is the text it was written as, with what `%20` stands for
// put back.
func dataURI(href string) ([]byte, error) {
	meta, body, ok := strings.Cut(strings.TrimPrefix(href, "data:"), ",")
	if !ok {
		return nil, fmt.Errorf("the data: address has no comma in it")
	}
	for _, part := range strings.Split(meta, ";") {
		if strings.EqualFold(strings.TrimSpace(part), "base64") {
			clean := strings.Map(func(r rune) rune {
				if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
					return -1
				}
				return r
			}, body)
			if raw, err := base64.StdEncoding.DecodeString(clean); err == nil {
				return raw, nil
			}
			return base64.RawStdEncoding.DecodeString(clean)
		}
	}
	s, err := url.PathUnescape(body)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// decodePicture is a PNG, a JPEG or a GIF turned into a canvas, at the size it
// was encoded at. The channels come out of the decoder premultiplied by their
// alpha and a canvas holds them straight, so the multiplication is undone the
// same way the icon reader in the android package does it: everything
// transparent would otherwise come out too dark to lay over anything.
func decodePicture(data []byte) (*canvas.Canvas, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("it is not a PNG, a JPEG or a GIF: %w", err)
	}
	b := src.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, fmt.Errorf("it has no size to it")
	}
	cv, err := canvas.NewCanvas(b.Dx(), b.Dy())
	if err != nil {
		return nil, err
	}
	for y := range b.Dy() {
		for x := range b.Dx() {
			r, g, bl, a := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			if a == 0 {
				cv.Put(x, y, canvas.Transparent)
				continue
			}
			cv.Put(x, y, canvas.RGBA(uint8(r*0xFF/a), uint8(g*0xFF/a), uint8(bl*0xFF/a), uint8(a>>8)))
		}
	}
	return cv, nil
}

// fadePicture is a picture at an opacity, in pixels of its own — the same
// multiply of alpha that [fade] does to a colour, done to every pixel of a
// bitmap instead. The picture it is given is left exactly as it is, since what
// is in the cache belongs to every element that points at the same address.
func fadePicture(src *canvas.Canvas, opacity float64) *canvas.Canvas {
	if opacity >= 1 {
		return src
	}
	cv, err := canvas.NewCanvas(src.Width, src.Height)
	if err != nil {
		// There is nowhere to put a faded copy, and the unfaded one is the
		// worse of the two things to paint: it is a canvas that could not be
		// allocated, which does not happen, and saying nothing beats painting
		// a hole where a picture was.
		return src
	}
	for y := range src.Height {
		for x := range src.Width {
			if c := src.At(x, y); c != 0 {
				cv.Put(x, y, fade(c, opacity))
			}
		}
	}
	return cv
}

// paintPicture lays an `<image>`'s picture into the box it was given, through
// the transform that puts the drawing on the canvas — the same way writing is
// painted, with the element's own transform composed on at paint time rather
// than put into anything while the drawing is read.
//
// The picture's own pixels first go to where they belong inside the box: the
// fit says how big they come out and the alignment says where they sit when
// that leaves anything over, which is [picturePlacement]. A picture that fits
// goes straight onto the canvas; a sliced one is cut to the box on a picture of
// its own first, since what is cut is the rectangle the drawing asked for and
// nothing beyond it, however the picture was turned on the way.
func paintPicture(cv *canvas.Canvas, n *Node, m canvas.Matrix) {
	p := n.Pic
	pic := p.Pixels
	if pic == nil || pic.Width <= 0 || pic.Height <= 0 {
		return
	}
	bw, bh := p.Box[2], p.Box[3]
	if bw <= 0 || bh <= 0 {
		return
	}
	sx, sy, ox, oy := picturePlacement(p, float64(pic.Width), float64(pic.Height))
	total := m.Mul(n.Style.Transform)
	matrix := total.Mul(canvas.Translate(ox, oy).Mul(canvas.Scale(sx, sy)))
	if p.Fit != PictureFitSlice {
		cv.BlitMatrix(pic, matrix, canvas.Area{Width: pic.Width, Height: pic.Height})
		return
	}
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		// There is no picture to cut: a canvas of no size has nothing to draw
		// on and nothing to draw into either, the same as every other layer
		// this package fails to make.
		return
	}
	layer.BlitMatrix(pic, matrix, sliceRegion(p.Box, ox, oy, sx, sy, pic))
	x0, y0, x1, y1 := regionArea(p.Box[0], p.Box[1], bw, bh, total)
	layer.MaskRect(x0, y0, x1-x0, y1-y0)
	cv.BlitOver(0, 0, layer)
}

// picturePlacement is where a picture's own pixels go inside its box: how much
// each axis of it is scaled by and which corner of the box it starts from. The
// fit decides the scale — the biggest that still fits for a meet, the biggest
// that covers for a slice, and one scale per axis for none — and the alignment
// decides what is left over of the box, a slice leaving a negative amount and
// pulling the picture past the edge it is not pinned to.
func picturePlacement(p *Picture, pw, ph float64) (sx, sy, ox, oy float64) {
	bx, by, bw, bh := p.Box[0], p.Box[1], p.Box[2], p.Box[3]
	if p.Fit == PictureFitNone {
		return bw / pw, bh / ph, bx, by
	}
	s := bw / pw
	if p.Fit == PictureFitSlice {
		if b := bh / ph; b > s {
			s = b
		}
	} else if b := bh / ph; b < s {
		s = b
	}
	fx, fy := p.Align.offsets()
	return s, s, bx + (bw-pw*s)*fx, by + (bh-ph*s)*fy
}

// sliceRegion is the stretch of the picture a sliced box can show: where the
// box falls in the picture's own pixels, which is the only part worth walking.
// The bounds are rounded outwards so that no pixel the box can reach is left
// out — the cut itself is the rectangle [paintPicture] takes to the layer,
// which is the exact one, and this is only about not walking a photograph from
// end to end to show a corner of it. The placement never turns (the turn is on
// the element above it), so the box lands in the picture as a rectangle however
// the drawing has moved.
func sliceRegion(box [4]float64, ox, oy, sx, sy float64, pic *canvas.Canvas) canvas.Area {
	x0 := max(0, int(math.Floor((box[0]-ox)/sx)))
	y0 := max(0, int(math.Floor((box[1]-oy)/sy)))
	x1 := min(pic.Width, int(math.Ceil((box[0]+box[2]-ox)/sx)))
	y1 := min(pic.Height, int(math.Ceil((box[1]+box[3]-oy)/sy)))
	if x1 <= x0 || y1 <= y0 {
		return canvas.Area{}
	}
	return canvas.Area{X: x0, Y: y0, Width: x1 - x0, Height: y1 - y0}
}
