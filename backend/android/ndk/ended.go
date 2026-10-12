//go:build android

package ndk

// Ended reports whether the codec has produced everything it is going to.
func (d *Decoder) Ended() bool { return d.ended }
