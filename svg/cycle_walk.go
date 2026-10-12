package svg

// cycleWalk is the state a mask/pattern cycle walk carries: the drawing it
// warns through, and which references are on the stack (gray) and which are
// finished (done).
type cycleWalk struct {
	img  *Image
	gray map[any]bool
	done map[any]bool
}

func (w *cycleWalk) visitMask(def *maskDef) {
	if w.done[def] {
		return
	}
	w.gray[def] = true
	w.walk(def.content)
	w.gray[def] = false
	w.done[def] = true
}

func (w *cycleWalk) visitPattern(p *pattern) {
	if w.done[p] {
		return
	}
	w.gray[p] = true
	w.walk(p.content)
	w.gray[p] = false
	w.done[p] = true
}
