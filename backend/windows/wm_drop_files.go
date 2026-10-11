//go:build windows

package windows

const (
	wmDropFiles = 0x0233

	cfUnicodeText = 13
	cfHDrop       = 15

	// What DragQueryFile answers with when it is asked for file 0xFFFFFFFF:
	// how many there are rather than the name of one.
	dragQueryCount = 0xFFFFFFFF
)
