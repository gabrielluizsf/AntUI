//go:build linux && !android

package antui

// pulsePlayback is what a create-playback-stream request needs to say.
type pulsePlayback struct {
	name    string
	audio   *Audio
	target  uint32
	minimum uint32
}
