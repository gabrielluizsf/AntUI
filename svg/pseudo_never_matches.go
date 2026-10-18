package svg

// neverMatches is the pseudo-classes a still drawing can never be in: the ones
// about user interaction, about the state of a form control, of a link, of a
// media element, or of the viewport. They are valid CSS, so a rule that names
// one is read and simply never matches, rather than being refused as though
// the name were unknown. The ones a drawing can answer — where an element
// stands, `:not`, `:is`, `:where`, `:has` — are read individually instead.
var neverMatches = map[string]bool{
	"hover": true, "active": true, "focus": true, "focus-visible": true,
	"focus-within": true, "checked": true, "disabled": true, "enabled": true,
	"indeterminate": true, "default": true, "valid": true, "invalid": true,
	"in-range": true, "out-of-range": true, "required": true, "optional": true,
	"read-only": true, "read-write": true, "placeholder-shown": true,
	"autofill": true, "user-invalid": true, "user-valid": true, "blank": true,
	"target": true, "target-within": true, "visited": true, "link": true,
	"any-link": true, "local-link": true, "playing": true, "paused": true,
	"current": true, "past": true, "future": true, "fullscreen": true,
	"modal": true, "picture-in-picture": true, "popover-open": true,
	"open": true, "closed": true, "defined": true, "host": true,
	"host-context": true, "state": true,
	"nth-col": true, "nth-last-col": true, "has-slotted": true,
	"active-view-transition": true,
}

// isPseudoElement says whether a name is a pseudo-element, which a still
// drawing draws no box for. The four legacy ones may be written with one colon
// and the rest with two.
func isPseudoElement(name string) bool {
	switch name {
	case "before", "after", "first-line", "first-letter", "marker",
		"backdrop", "selection", "placeholder", "cue", "file-selector-button",
		"grammar-error", "spelling-error", "view-transition", "details-content",
		"target-text", "highlight", "part", "slotted", "cue-region":
		return true
	}
	return false
}
