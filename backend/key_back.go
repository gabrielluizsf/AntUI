package backend

// Keys a phone has and a desktop does not. KeyBack is the one that matters:
// it is the system's "go back", and a program that does not handle it gets
// the platform's answer, which is to close.
const (
	KeyBack Key = 256 + 35 + iota
	KeyMenu
	KeySearch
	KeyVolumeUp
	KeyVolumeDown
)
