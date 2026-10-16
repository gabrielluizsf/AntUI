package svg

// clipNow is the clip one node inside a `<clipPath>` takes while the shapes
// of the clipPath it is inside are being worked out: the shapes of the
// clipPath it names, moved into the node's own place the same way
// [Image.resolveClip] moves them. A name that points back at a clipPath
// whose shapes are being worked out right now is a circle, and it is left
// off with one warning saying so — a clip left off keeps everything, so the
// shape it was written on stands whole inside the clipPath, which is what
// lets the walk come back rather than going round for ever.
func (img *Image) clipNow(n *Node, id string, gray map[string]bool) *clipRegion {
	warn := func(format string, args ...any) { img.warnNamed(n.Name, format, args...) }
	if gray[id] {
		warn("the clip-path here points back at the <clipPath> this picture is inside, which would go round for ever, so the clip is left off")
		return nil
	}
	cp := img.measureClip(id, gray)
	if cp == nil {
		// Not there: resolveClip would have said so and kept nothing while
		// the picture was built, so a name only reaches here from a clipPath
		// that was standing when it was read.
		return &clipRegion{}
	}
	return clipUnder(img.clipRegionOf(cp, n.Style), n.Style, n, warn)
}
