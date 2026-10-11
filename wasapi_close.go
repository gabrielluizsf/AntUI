//go:build windows

package antui

func (w *wasapi) close() error {
	select {
	case <-w.quit:
		return nil // already asked
	default:
	}
	close(w.quit)

	// Wake the thread so it notices, rather than leaving it asleep until the
	// device happens to want another block.
	if w.event != 0 {
		procSetEvent.Call(w.event)
	}
	<-w.done
	return nil
}
