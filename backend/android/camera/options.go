//go:build android

package camera

// Options is how a preview is opened.
type Options struct {
	// Width and Height are what to ask the camera for. They are a request:
	// the camera picks the nearest size it supports, and the frames that
	// come back say what it chose. Zero asks for 1280x720, which every
	// camera has.
	Width, Height int
	// Upright turns each frame so that it matches what the user is looking
	// at, which is what almost every app wants and what nothing does by
	// itself. Turning it off gives the sensor's own image, sideways.
	Upright bool
}
