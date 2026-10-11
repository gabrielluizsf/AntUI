package svg

import "strings"

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
