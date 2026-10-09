package css

// Transition is one property an element animates between its computed values
// when the cascade changes it: the transition shorthand and its longhands
// come together here. Prop is the canonical property name (such as "opacity"
// or "background-color") or "all"; Duration and Delay are the two times CSS
// writes in either order; Timing is the easing. A zero Duration means the
// property snaps, exactly as a style without a transition does.
type Transition struct {
	Prop     string
	Duration Time
	Timing   Timing
	Delay    Time
}
