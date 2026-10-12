//go:build android

package ndk

// Video reports whether this is a picture track.
func (t Track) Video() bool { return len(t.Mime) > 6 && t.Mime[:6] == "video/" }
