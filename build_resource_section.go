package antui

func buildResourceSection(pictures []iconPicture) resourceSection {
	blobs := resourceBlobs(pictures)
	l := layoutResourceSection(blobs)
	writeResourceSection(l, blobs)
	return resourceSection{data: l.out, fixups: l.leaves}
}
