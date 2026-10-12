//go:build android

package speech

// Say speaks, and returns as soon as the words are queued.
func (s *Speaker) Say(text string, mode int32) error {
	_, err := s.say(text, mode, false)
	return err
}
