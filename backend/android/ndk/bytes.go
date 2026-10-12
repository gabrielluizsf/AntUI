//go:build android

package ndk

// Bytes is how much one pixel takes.
func (f Format) Bytes() int {
	switch f {
	case RGB565:
		return 2
	case RGB888:
		return 3
	}
	return 4
}
