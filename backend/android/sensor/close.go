//go:build android

package sensor

// Close stops every sensor and gives the thread back. It waits for the
// thread to finish, so that nothing is still writing when it returns.
func (s *Stream) Close() error {
	s.stopOnce.Do(func() {
		close(s.stop)
		<-s.stopped
	})
	return nil
}
