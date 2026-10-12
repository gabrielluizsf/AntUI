//go:build android

package speech

// SayAndWait speaks and waits until the words have been said.
//
// **Not from the frame loop**: it takes as long as the sentence does.
func (s *Speaker) SayAndWait(text string, mode int32) error {
	done, err := s.say(text, mode, true)
	if err != nil {
		return err
	}
	return <-done
}
