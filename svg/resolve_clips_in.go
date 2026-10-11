package svg

// resolveClipsIn follows every `clip-path` still waiting in a picture: they
// were named while the picture was built, when the clipPath they point at
// had a name but no shapes yet, and now that it has, the shapes are moved
// into the place of whoever named them the same way [Image.resolveClip] does
// for a reference anywhere else in the drawing.
func (img *Image) resolveClipsIn(n *Node, gray map[string]bool) {
	if n == nil {
		return
	}
	if n.clip != nil && n.clip.pending != "" {
		n.clip = img.clipNow(n, n.clip.pending, gray)
	}
	for _, k := range n.Kids {
		img.resolveClipsIn(k, gray)
	}
}
