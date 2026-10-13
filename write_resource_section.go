package antui

import "encoding/binary"

// writeResourceSection fills a laid-out section: the directory tree and the
// leaves.
func writeResourceSection(l resourceLayout, blobs []resourceBlob) {
	put32 := func(offset, value uint32) {
		binary.LittleEndian.PutUint32(l.out[offset:], value)
	}
	put16 := func(offset uint32, value uint16) {
		binary.LittleEndian.PutUint16(l.out[offset:], value)
	}

	// A directory header: no name entries, and however many by number.
	dir := func(offset uint32, count int) {
		put16(offset+12, 0)             // named entries
		put16(offset+14, uint16(count)) // entries by number
	}
	// One entry of a directory: what it is called, and where it points.
	entry := func(offset uint32, id uint16, target uint32, isDir bool) {
		put32(offset, uint32(id))
		if isDir {
			target |= resourceSubdir
		}
		put32(offset+4, target)
	}

	// The root, by type.
	dir(0, 2)
	entry(resourceDirSize, rtIcon, l.typeDirs[rtIcon], true)
	entry(resourceDirSize+resourceEntrySize, rtGroupIcon, l.typeDirs[rtGroupIcon], true)

	// The icons, by number.
	dir(l.typeDirs[rtIcon], l.icons)
	for i := range l.icons {
		entry(l.typeDirs[rtIcon]+resourceDirSize+uint32(i)*resourceEntrySize,
			uint16(i+1), l.nameDirs[i], true)
		// And its one language.
		dir(l.nameDirs[i], 1)
		entry(l.nameDirs[i]+resourceDirSize, langNeutral, l.leaves[i], false)
	}

	// The group, which is the one Windows looks up to find the rest.
	dir(l.typeDirs[rtGroupIcon], 1)
	entry(l.typeDirs[rtGroupIcon]+resourceDirSize, 1, l.groupNameDir, true)
	dir(l.groupNameDir, 1)
	entry(l.groupNameDir+resourceDirSize, langNeutral, l.leaves[len(blobs)-1], false)

	// The leaves, and the bytes they point at.
	for i, one := range blobs {
		put32(l.leaves[i], l.starts[i]) // the linker adds the section's address
		put32(l.leaves[i]+4, uint32(len(one.data)))
		put32(l.leaves[i]+8, 0) // code page: none, the data is not text
		put32(l.leaves[i]+12, 0)
		copy(l.out[l.starts[i]:], one.data)
	}
}
