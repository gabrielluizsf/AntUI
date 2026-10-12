//go:build android

package gallery

// collection is the address of the library, and folder the directory inside
// shared storage it lives in.
func (k Kind) collection() (class, field string) {
	switch k {
	case Video:
		return "android/provider/MediaStore$Video$Media", "EXTERNAL_CONTENT_URI"
	case Audio:
		return "android/provider/MediaStore$Audio$Media", "EXTERNAL_CONTENT_URI"
	}
	return "android/provider/MediaStore$Images$Media", "EXTERNAL_CONTENT_URI"
}
