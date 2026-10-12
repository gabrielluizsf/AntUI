//go:build android

package gallery

// String names the library, which is what an error message wants.
func (k Kind) String() string {
	switch k {
	case Video:
		return "video"
	case Audio:
		return "audio"
	}
	return "image"
}
