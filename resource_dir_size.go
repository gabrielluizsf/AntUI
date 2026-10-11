package antui

const (
	resourceDirSize   = 16
	resourceEntrySize = 8
	resourceLeafSize  = 16
	// The high bit of an entry's offset says it points at another directory
	// rather than at a leaf.
	resourceSubdir = 0x80000000
)
