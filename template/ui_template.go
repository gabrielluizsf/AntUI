package template

import "github.com/gabrielluizsf/antui/template/event"

// uiTemplate is [New]'s ordinary answer to [Template]: it runs the shared
// interaction machinery and leaves the painting to the style and the hearing
// to the bank.
type uiTemplate struct {
	style            Style
	sound            event.SoundBank
	sliderDragOffset int
	sliderDragID     uint64
	selectOpen       uint32 // the widget id of an open dropdown menu, 0 when none
	selectCursor     int    // which option the keyboard is walking in an open menu
	dateOpen         uint32 // the widget id of an open calendar popup, 0 when none
	dateYear         int    // the month an open calendar popup shows
	dateMonth        int
	dateCursor       int  // the day an open calendar popup highlights
	datePickMonth    bool // whether that popup is showing the month grid instead of the day grid

	// ld is the start-up loading screen this template wears until the
	// program says Ready — see [loading]. It is nil for a template built
	// as a bare struct instead of through [New], and a nil screen is one
	// that is never active, so such a template simply has no screen.
	ld *loading
}
