//go:build android

package ndk

import "time"

// Track is one stream inside a media file: a picture, a sound, a subtitle.
type Track struct {
	Index int
	// Mime is what it holds — "video/avc", "audio/mp4a-latm".
	Mime string
	// Duration is how long the whole file is, which the container reports
	// per track and which is usually the same for all of them.
	Duration time.Duration

	// Width and Height are set on a video track, and Rotation is how far it
	// has to be turned to look right — a video shot on a phone held sideways
	// carries 90 here and is stored as if it were not.
	Width, Height int
	Rotation      int
	FrameRate     int

	// SampleRate and Channels are set on an audio track.
	SampleRate int
	Channels   int
}
