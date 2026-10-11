package canvas

import "sync"

var (
	builtinMu     sync.Mutex
	builtinScaled = map[int]*Face{}
	// builtinSized is the built-in font at sizes asked for outright rather
	// than as a multiple, keyed by the whole pixel it was made at.
	builtinSized = map[int]*Face{}
)
