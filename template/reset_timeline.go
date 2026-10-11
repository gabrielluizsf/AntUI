package template

// resetTimeline opens a new frame: the widget ordinals restart at zero and
// the current entry reference clears, so the next layout begins a fresh
// run. An entry not drawn for two frames is dropped, so a widget that
// stopped appearing (or changed role and header) does not linger in the map.
func (cs *cssStyle) resetTimeline() {
	cs.frame++
	cs.seq = 0
	cs.cur = nil
	for k, e := range cs.timeline {
		if e.seen+2 <= cs.frame {
			delete(cs.timeline, k)
		}
	}
}
