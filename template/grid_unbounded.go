package template

import "math"

// gridUnbounded is the room handed to a track that is being asked how wide it
// wants to be: no container outside has an opinion yet, so its ceiling is the
// whole int and only the track's own floor and ceiling say no.
const gridUnbounded = math.MaxInt
