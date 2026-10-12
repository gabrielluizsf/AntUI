//go:build android

package ndk

// Audio reports whether this is a sound track.
func (t Track) Audio() bool { return len(t.Mime) > 6 && t.Mime[:6] == "audio/" }
