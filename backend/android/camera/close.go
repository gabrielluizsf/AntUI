//go:build android

package camera

// Close stops the camera and gives it back. A camera left open is a camera
// no other app can use, and the platform does not take it back on its own.
func (p *Preview) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	p.device.Close()
	p.manager.Close()
	p.canvas = nil
	return nil
}
