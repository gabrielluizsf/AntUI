//go:build darwin && cgo

package antui

// fillQueue is the Go half of the callback, with the C types left at the door
// so that the part which can be got wrong is the part a test can call.
func fillQueue(id int, block []float32) {
	queuesMu.Lock()
	c := queues[id]
	queuesMu.Unlock()

	if c == nil {
		// The stream was closed between the queue asking and this answering.
		for i := range block {
			block[i] = 0
		}
		return
	}
	c.audio.pull(block)
}
