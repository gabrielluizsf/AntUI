package template

import "strconv"

// beginWidget opens the timeline of the widget being laid out. The key is
// the role plus its ordinal in this frame plus its label, which pins a
// widget through its frame and back to the same entry next frame. The entry
// becomes cs.cur, so the painter's styleFor calls land on it.
func (cs *cssStyle) beginWidget(role, label string) *widgetEntry {
	cs.seq++
	key := role + "\x00" + strconv.Itoa(cs.seq) + "\x00" + label
	e := cs.timeline[key]
	if e == nil {
		e = &widgetEntry{role: role, fresh: true}
		cs.timeline[key] = e
	}
	e.seen = cs.frame
	cs.cur = e
	return e
}
