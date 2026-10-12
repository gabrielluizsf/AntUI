package svg

// textPathRuns reads the writing a `<textPath>` hangs along the outline its href
// names, with `startOffset` saying how far along it the first letter goes. What
// it names could not be followed, and then there is nowhere to write, which is
// worth one sentence rather than the writing turning up in the wrong place.
func (img *Image) textPathRuns(e *element, st Style, warn func(string, ...any)) []TextRun {
	along := img.textPathOutline(e, warn)
	if along == nil {
		return nil
	}
	kid := st.with(e, warn)
	img.resolvePaints(&kid, warn)
	sub := img.textRuns(e, kid, warn)
	off := startOffset(e, along, warn)
	for i := range sub {
		sub[i].Path, sub[i].Offset = along, off
	}
	return sub
}
