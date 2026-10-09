package css

// DefaultFontSize is the reference-pixel font the engine assumes when nothing
// sets font-size; rem and the font-relative units use it.
const DefaultFontSize = 16

// Units is the measurement context a Length needs to become pixels: the
// containing window (vw/vh/vmin/vmax), the containing width (percentages) and
// the font sizes (em/ch/ex, rem). Zero fields fall back to sensible defaults.
type Units struct {
	Width, Height int // the window, for viewport units and %-of-width
	Font          int // the element's font size, for em/ch/ex
	Root          int // the document font size, for rem

	// Scale and Scheme are the window's display density and the colour
	// scheme its system paints in — what the resolution and
	// prefers-color-scheme media features read. Zero and SchemeUnknown are
	// "the system did not say".
	Scale  float64
	Scheme Scheme
}

// base is the fallback font when a context has not said anything.
func (u Units) font() int { return max(u.Font, DefaultFontSize) }

func (u Units) root() int { return max(u.Root, DefaultFontSize) }
