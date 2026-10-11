package antui

// layoutResourceSection places the directories first, then the leaves, then
// the bytes, each on an eight-byte boundary.
func layoutResourceSection(blobs []resourceBlob) resourceLayout {
	icons := len(blobs) - 1

	at := uint32(resourceDirSize + 2*resourceEntrySize) // the root: two types
	typeDirs := map[uint16]uint32{}

	typeDirs[rtIcon] = at
	at += uint32(resourceDirSize + icons*resourceEntrySize)
	nameDirs := make([]uint32, 0, icons+1)
	for range icons {
		nameDirs = append(nameDirs, at)
		at += resourceDirSize + resourceEntrySize // one language
	}

	typeDirs[rtGroupIcon] = at
	at += resourceDirSize + resourceEntrySize
	groupNameDir := at
	at += resourceDirSize + resourceEntrySize

	leaves := make([]uint32, len(blobs))
	for i := range blobs {
		leaves[i] = at
		at += resourceLeafSize
	}

	starts := make([]uint32, len(blobs))
	for i := range blobs {
		at = (at + 7) &^ 7
		starts[i] = at
		at += uint32(len(blobs[i].data))
	}

	return resourceLayout{
		out:          make([]byte, at),
		typeDirs:     typeDirs,
		nameDirs:     nameDirs,
		groupNameDir: groupNameDir,
		leaves:       leaves,
		starts:       starts,
		icons:        icons,
	}
}
