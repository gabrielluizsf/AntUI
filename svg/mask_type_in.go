package svg

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// maskTypeIn is how one `<mask>` measures how much it keeps: `luminance`,
// which is what a mask means where nothing has said otherwise, or `alpha`,
// which keeps by how much of the mask is there rather than by how bright it
// is. It is read from the attribute and from a declaration in the `style`,
// the way clip-rule is, and it is not inherited — a mask is measured by what
// it says about itself and not by the place it stands in.
//
// A value that is neither is said out loud and read as the default, because
// the two measures disagree about the whole of a black mask and an author who
// asked for one of them should not be given the other in silence.
func maskTypeIn(e *element, warn func(string, ...any)) canvas.MaskMode {
	raw, ok := e.attr("mask-type"), e.hasAttr("mask-type")
	if !ok {
		if s := e.attr("style"); s != "" {
			raw, ok = parseDeclarations(s)["mask-type"]
		}
	}
	if !ok {
		return canvas.MaskLuminance
	}
	switch keyword(raw) {
	case "alpha":
		return canvas.MaskAlpha
	case "luminance", "":
		return canvas.MaskLuminance
	default:
		warn("the mask-type %q is not luminance or alpha, so the mask measures how bright it is", raw)
		return canvas.MaskLuminance
	}
}
