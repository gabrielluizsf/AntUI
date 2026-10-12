//go:build android

package audio

// Rate is how many frames a second the device is really recording at.
func (r *Recorder) Rate() int { return r.stream.Rate }
