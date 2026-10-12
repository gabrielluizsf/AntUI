//go:build android

package sensor

func (s *Stream) deliver(r Reading) {
	s.mu.Lock()
	s.latest[r.Kind] = r
	s.fresh[r.Kind] = true
	s.mu.Unlock()
	if s.events != nil {
		select {
		case s.events <- r:
		default:
			// The buffer is full and nobody is reading. Dropping the newest
			// would be wrong — for a sensor the newest is the one that
			// matters — so the oldest goes.
			select {
			case <-s.events:
			default:
			}
			select {
			case s.events <- r:
			default:
			}
		}
	}
}
