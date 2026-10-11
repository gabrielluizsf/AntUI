//go:build linux && !android

package antui

// openAudio is the Linux half of OpenAudio.
func openAudio(name string, a *Audio) (driver, error) {
	p, err := dialPulse(name, a)
	if err != nil {
		return nil, err
	}
	if err := p.startStream(name, a); err != nil {
		p.conn.Close()
		return nil, err
	}
	return p, nil
}
