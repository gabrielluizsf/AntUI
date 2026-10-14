//go:build android

package gallery

// Entry describes one item being added to the device's library.
type Entry struct {
	Kind  Kind
	Name  string
	Mime  string
	Album string
}
