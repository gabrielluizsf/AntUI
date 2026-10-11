package antui

// trackVelocity smooths the finger's speed over frames. One frame's movement
// is far too noisy to throw anything with — a finger that stalls for a single
// frame before lifting would otherwise register as a dead stop.
func (win *Window) trackVelocity(now float64, t *Touch) {
	g := &win.gesture
	if dt := now - g.lastTime; dt > 0 {
		const smoothing = 0.6
		vx := float64(t.X-g.lastX) / dt
		vy := float64(t.Y-g.lastY) / dt
		g.vx = g.vx*(1-smoothing) + vx*smoothing
		g.vy = g.vy*(1-smoothing) + vy*smoothing
		g.lastX, g.lastY, g.lastTime = t.X, t.Y, now
	}
}
