//go:build android

package audio

// Err is what stopped it, or nil.
func (r *Recorder) Err() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.err
}
