package antui

// resourceBlobs is the leaves of the section: one per icon, then one for the
// group that ties them together.
func resourceBlobs(pictures []iconPicture) []resourceBlob {
	blobs := make([]resourceBlob, 0, len(pictures)+1)
	for i, picture := range pictures {
		blobs = append(blobs, resourceBlob{uint16(i + 1), rtIcon, picture.data})
	}
	return append(blobs, resourceBlob{1, rtGroupIcon, groupIcon(pictures)})
}
