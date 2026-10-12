//go:build android

package sensor

// Events is every reading, for the sensors that report an event rather than
// a value — a step, a significant motion. It is nil unless [Options.Buffer]
// was set.
func (s *Stream) Events() <-chan Reading { return s.events }
