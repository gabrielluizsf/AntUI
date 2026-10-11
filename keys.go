package antui

import (
	"iter"
)

// Keys walks every key that is held or that went down or up this frame.
func (win *Window) Keys() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		if win == nil {
			return
		}
		for k := range keyCount {
			if !win.keyState[k] && !win.keyDownFrame[k] && !win.keyUpFrame[k] {
				continue
			}
			if !yield(Key(k)) {
				return
			}
		}
	}
}
