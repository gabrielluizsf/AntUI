//go:build windows

package antui

// openAudio is the Windows half of OpenAudio.
//
// Everything happens on the stream's own thread, because COM is threaded: an
// interface pointer belongs to the apartment it was made in, and using it
// from another goroutine is undefined in a way that works right up until it
// does not.
func openAudio(name string, a *Audio) (driver, error) {
	w := &wasapi{
		audio: a,
		ready: make(chan error, 1),
		quit:  make(chan struct{}),
		done:  make(chan struct{}),
	}
	go w.run()
	if err := <-w.ready; err != nil {
		return nil, err
	}
	return w, nil
}
