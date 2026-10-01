package template

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/svg"
)

// The icons the built-in widgets draw themselves, by the name they are
// registered under. A widget that draws one of these asks for it by name, so a
// developer who replaces it changes what every widget of that kind looks like
// at once.
const (
	// IconCheck is the tick inside a checked checkbox.
	IconCheck = "check"
	// IconDot is the filled centre of a selected radio button.
	IconDot = "dot"
	// IconCaret is the arrow that says a dropdown opens downwards.
	IconCaret = "caret"
	// IconCaretUp is the same arrow the other way up, which is what a date
	// picker shows while its calendar is open.
	IconCaretUp = "caret-up"
)

// The built-in icons, written the way an icon set is: a twenty-four unit square
// of a drawing, stroked rather than filled so the line weight stays even at any
// size, with rounded ends and corners so a twelve-pixel tick does not look
// chipped, and painted in `currentColor` so one drawing serves every colour a
// widget wants to say it in.
//
// The stroke width is why a tick drawn here does not look like the four straight
// lines it replaces: a stroked path is one line of a given weight, so it is
// equally thick along the diagonal and at the corner, which is what a tick is.
var builtinIcons = map[string]string{
	IconCheck: `<svg viewBox="0 0 24 24">
		<path d="M5 12.5 L9.5 17 L19 7" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
	IconDot: `<svg viewBox="0 0 24 24">
		<circle cx="12" cy="12" r="6" fill="currentColor"/>
	</svg>`,
	IconCaret: `<svg viewBox="0 0 24 24">
		<path d="M5 9.5 L12 16 L19 9.5" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
	IconCaretUp: `<svg viewBox="0 0 24 24">
		<path d="M5 14.5 L12 8 L19 14.5" fill="none" stroke="currentColor"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>`,
}

// An icon as it is registered: the text it was written in, so a developer can
// read it and change it, and the drawing that text became, read once so that
// painting it is not parsing it every frame.
type iconEntry struct {
	src string
	img *svg.Image
}

// The registered icons, and the pictures already painted from them.
//
// The lock is because a window paints on its own goroutine while the program
// that owns it may replace an icon at any moment; the map is read on every
// widget of every frame, so the read side is the one that matters and takes the
// lock only for as long as it takes to read one pointer.
var (
	iconMu sync.RWMutex
	icons  = map[string]*iconEntry{}

	// Painted icons, keyed by what they were painted at. A drawing remembers one
	// size and colour of itself, which is not enough for an interface that holds
	// a list of buttons: two of them in the same state share a picture, but the
	// same button hovered and not hovered alternates between two colours every
	// frame, and a drawing that could only remember one of them would redraw the
	// icon on every frame forever. So the pictures are kept here instead, where a
	// key covers all three.
	iconPaintMu sync.RWMutex
	iconPaint   = map[iconPaintKey]*canvas.Canvas{}
)

// iconPaintKey is one picture: which icon, how big, in which colour.
type iconPaintKey struct {
	name string
	w, h int
	tint canvas.Color
}

// iconPaintLimit is how many painted icons are kept before the oldest are thrown
// out by starting over. An interface has a handful of icons at a handful of
// sizes in a handful of states, so the map is only ever as big as the set of
// things being drawn — and the bound is what keeps a program that names icons
// by something that changes per frame from growing it without end.
const iconPaintLimit = 64

func init() { resetIcons() }

// resetIcons puts the built-in icons back, which is what a program wants after
// it has been changing them, and what a test wants between cases.
func resetIcons() {
	iconMu.Lock()
	defer iconMu.Unlock()
	icons = make(map[string]*iconEntry, len(builtinIcons))
	for name, src := range builtinIcons {
		// A built-in icon is known to be readable, since it is the one this
		// package ships; one that is not would be a mistake here rather than
		// something to report, and it would draw as nothing.
		img, err := svg.Parse(src)
		if err != nil {
			continue
		}
		icons[name] = &iconEntry{src: src, img: img}
	}
	iconPaintMu.Lock()
	iconPaint = map[iconPaintKey]*canvas.Canvas{}
	iconPaintMu.Unlock()
}

// SetIcon registers the drawing a named icon is made of, replacing whatever was
// there, and reports a drawing that could not be read — in which case the icon
// that was registered stays as it was.
//
// The drawing is SVG text, which is the point: a developer reads the icon that
// is already there with [IconSource], changes it, and hands the changed text
// back. Recolouring it, thickening the line, squaring off the corner of a caret
// or drawing a label of one's own across the tick of a checkbox are all changes
// to the text, and none of them need anything else from the engine.
//
//	btn := template.IconSource(template.IconCheck)
//	btn = strings.Replace(btn, "stroke-width", `fill="red" stroke-width`)
//	template.SetIcon(template.IconCheck, btn)
//
// The name is what a stylesheet's `url()` and the widgets ask for, so a new
// name makes a new icon: `SetIcon("save", src)` is what a button wears when its
// stylesheet says `background-image: url(save)`.
func SetIcon(name, src string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("template: an icon needs a name")
	}
	img, err := svg.Parse(src)
	if err != nil {
		return fmt.Errorf("template: the icon %q could not be read: %w", name, err)
	}
	iconMu.Lock()
	icons[name] = &iconEntry{src: src, img: img}
	iconMu.Unlock()
	// The pictures already painted are of the icon that was there before, so they
	// go: what is on screen for a frame is the old one, which is no worse than a
	// frame painted before the icon was replaced.
	iconPaintMu.Lock()
	for key := range iconPaint {
		if key.name == name {
			delete(iconPaint, key)
		}
	}
	iconPaintMu.Unlock()
	return nil
}

// IconSource is the SVG text a named icon was registered from, which is the
// drawing to start from when changing it. The second answer says whether there
// is an icon of that name at all.
func IconSource(name string) (string, bool) {
	iconMu.RLock()
	defer iconMu.RUnlock()
	e := icons[name]
	if e == nil {
		return "", false
	}
	return e.src, true
}

// IconNames are the icons registered right now, sorted, which is what a tool
// that lets a developer pick an icon wants to show.
func IconNames() []string {
	iconMu.RLock()
	defer iconMu.RUnlock()
	names := make([]string, 0, len(icons))
	for name := range icons {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// IconSize is the size an icon is drawn at by default, in pixels, for the face
// of a widget that does not scale. An icon is asked for by the size it should
// be, so this is only the number a caller uses when it has no scale of its own.
const IconSize = 16

// DrawIcon paints a named icon into the rectangle given, in the colour given,
// and says whether there was one to paint.
//
// The icon is fitted into the rectangle keeping its own proportions and centred
// in what is left over, so a square rectangle gives an icon drawn square and a
// rectangle that is not square gives the same icon with its margins, which is
// what every icon set does and what keeps a tick from being stretched into a
// different mark.
//
// Painting nothing is not an error: an icon that is not registered simply does
// not appear, so a widget that asks for one it may not have still draws.
func DrawIcon(cv *canvas.Canvas, name string, x, y, w, h int, tint canvas.Color) bool {
	if cv == nil {
		return false
	}
	painted := IconCanvas(name, w, h, tint)
	if painted == nil {
		return false
	}
	cv.BlitScaled(x, y, w, h, painted)
	return true
}

// IconCanvas is the picture a named icon makes at a size in a colour, already
// painted and ready to be blitted, or nil when there is no such icon. It is
// what [DrawIcon] draws and what a stylesheet's `url()` puts on a box, so both
// ways of showing an icon share one picture and one cache.
//
// The canvas is the icon itself and belongs to the registry: blit from it as
// many times as a frame needs and do not draw on it.
func IconCanvas(name string, w, h int, tint canvas.Color) *canvas.Canvas {
	if w <= 0 || h <= 0 {
		return nil
	}
	iconMu.RLock()
	e := icons[name]
	iconMu.RUnlock()
	if e == nil || e.img == nil {
		return nil
	}
	key := iconPaintKey{name: name, w: w, h: h, tint: tint}
	iconPaintMu.RLock()
	painted, ok := iconPaint[key]
	iconPaintMu.RUnlock()
	if ok {
		return painted
	}
	painted = e.img.RenderWith(w, h, tint)
	if painted == nil {
		return nil
	}
	iconPaintMu.Lock()
	// Starting over rather than keeping the oldest: an interface asks for the
	// same few pictures every frame, so what is worth keeping is decided by this
	// frame and not by what happened to be asked for first.
	if len(iconPaint) >= iconPaintLimit {
		iconPaint = make(map[iconPaintKey]*canvas.Canvas, iconPaintLimit)
	}
	iconPaint[key] = painted
	iconPaintMu.Unlock()
	return painted
}

// DrawIconIn is [DrawIcon] at the size a widget draws its own marks at, centred
// in the box given. It is what the built-in widgets call, so that the mark they
// draw is the same size and in the same place whatever the scale.
func DrawIconIn(cv *canvas.Canvas, name string, x, y, size int, tint canvas.Color) bool {
	return DrawIcon(cv, name, x, y, size, size, tint)
}
