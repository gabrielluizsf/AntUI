//go:build android

package speech

// How a new piece of speech meets one already being said.
const (
	// Replace stops whatever is being said and says this instead. It is what
	// almost everything wants: the thing being said is out of date the
	// moment there is something newer to say.
	Replace int32 = 0
	// Queue waits its turn.
	Queue int32 = 1
)
