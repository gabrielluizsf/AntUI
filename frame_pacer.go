package antui

import (
	"time"
)

// framePacer is the optional half of a backend, and it is told when the next
// frame is due rather than how long a frame may take: a platform that has to
// wait for something to happen can spend the time until then asleep in the
// kernel instead of coming straight back to be called again — the difference
// between an idle window costing nothing and an idle window costing a core.
// One that cannot wait is unaffected, since the core does the waiting itself.
type framePacer interface {
	setFrameDeadline(due time.Time)
}
