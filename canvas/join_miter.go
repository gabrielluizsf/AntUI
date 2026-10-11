package canvas

// A miter join cuts the outside of the corner off square, a round one rounds it
// off, and a bevel one chops it flat. A miter longer than its limit falls back
// to the bevel, because a corner drawn out to a spike many times its own width
// is not what "sharp" was meant to say.
const (
	JoinMiter LineJoin = iota
	JoinRound
	JoinBevel
)
