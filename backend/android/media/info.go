//go:build android

package media

import (
	"time"
)

// Info is what a file holds.
type Info struct {
	Duration time.Duration
	// Width and Height are the picture's, after the rotation the file
	// carries has been applied — which is the size it should be drawn at.
	Width, Height int
	// Rotation is how far the picture is stored turned. A video shot on a
	// phone held sideways is stored as if it were not, with a 90 here.
	Rotation int
	// SampleRate and Channels are the sound's, and are zero when the file
	// has none.
	SampleRate, Channels int
	// HasVideo and HasAudio say which tracks are there.
	HasVideo, HasAudio bool
}
