package css

import (
	"sync"
)

// Sheet is a parsed stylesheet: every rule in order, the @keyframes blocks it
// named, plus warnings for the declarations it did not understand. Like a
// browser, the engine skips what it does not know and keeps going; the
// warnings let a template tell its author.
type Sheet struct {
	rules     []*Rule
	order     int
	keyframes map[string]*Keyframes
	fonts     []*FontFace

	// fontCache holds what a Font query answered — itself, or nil when the
	// sheet holds none of its names — so a frame walks the list once rather
	// than every line it draws. It is thrown away when a @font-face is
	// added, which happens only while the sheet is still being read.
	fontCache map[fontKey]*FontFace
	fontMu    sync.RWMutex

	Warn []string
}
