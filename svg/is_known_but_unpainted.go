package svg

// isKnownButUnpainted says whether a tag is one SVG defines that paints nothing
// where it stands. They are left out quietly: a `<defs>` or a `<title>` in the
// middle of an icon is normal, and warning about each of them would bury the
// warnings that matter.
//
// A tag that paints where it stands is not in here: `<image>` draws a picture
// it fetches itself — see [Image.imageNode] — and it is built whole when its
// turn comes, the way writing is rather than the way a shape is. A `<symbol>`
// is the other way round — it draws nothing where it stands and everything it
// holds is drawn where a `<use>` points at it. A `<text>` is neither: it paints
// where it stands, and is built as writing rather than as a shape.
func isKnownButUnpainted(name string) bool {
	switch name {
	case "defs", "title", "desc", "metadata", "style", "styleSheet", "script",
		"view", "font", "fontFace", "clipPath", "mask", "pattern", "symbol",
		"linearGradient", "radialGradient", "stop", "filter", "marker", "set",
		"animate", "animateTransform", "animateMotion":
		return true
	}
	return false
}
