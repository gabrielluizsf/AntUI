package css

// Scheme is what the color-scheme media feature is read against: what the
// system paints its own interface in. SchemeUnknown is a system that did not
// say, and a query over it is treated as satisfied — the same answer any
// condition the canvas cannot measure gets, so a rule is never dropped for
// want of an answer.
type Scheme uint8

// The three answers a system can give.
const (
	SchemeUnknown Scheme = iota
	SchemeLight
	SchemeDark
)
