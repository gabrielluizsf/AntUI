//go:build android

package app

// String names the orientation the way Android does.
func (o Orientation) String() string {
	switch o {
	case Portrait:
		return "portrait"
	case Landscape:
		return "landscape"
	case Square:
		return "square"
	}
	return "unknown"
}
