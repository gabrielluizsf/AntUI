package antui

// trackFinger picks the finger recognizeOne watches, starting tracking on a
// newly-landed one, and reports nil when there is none to watch.
func (win *Window) trackFinger(now float64) *Touch {
	g := &win.gesture
	touches := win.touches
	if len(touches) == 0 {
		g.tracking = false
		return nil
	}

	// The finger being tracked, or the first one if there is no tracking yet
	// or the tracked one has gone.
	i := -1
	if g.tracking {
		i = win.touchIndex(g.id)
	}
	if i < 0 {
		t := touches[0]
		if !t.Began {
			// A finger that was already down when tracking was lost — after
			// a second finger came and went, say. Picking it up mid-gesture
			// would report a tap the user never made.
			g.tracking = false
			return nil
		}
		*g = gestureState{
			tracking:  true,
			id:        t.ID,
			startX:    t.X,
			startY:    t.Y,
			startTime: now,
			lastX:     t.X,
			lastY:     t.Y,
			lastTime:  now,
			tapTime:   g.tapTime,
			tapX:      g.tapX,
			tapY:      g.tapY,
			tapCount:  g.tapCount,
		}
		i = 0
	}
	return &touches[i]
}
