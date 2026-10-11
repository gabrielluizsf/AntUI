package antui

// an iconPicture is one size, as the bytes that go in the file.
type iconPicture struct {
	width, height int
	data          []byte // a PNG, which every Windows since Vista reads
}
