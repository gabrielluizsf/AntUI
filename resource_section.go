package antui

// resourceSection is the bytes of the section and where the addresses in it
// are, so that relocations can be written for them.
type resourceSection struct {
	data []byte
	// fixups are offsets into data holding a section-relative address the
	// linker has to turn into a real one.
	fixups []uint32
}
