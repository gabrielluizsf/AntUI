package css

// motionStyle carries the animation fields of a Style. Transitions is the
// finished transition list: which properties animate toward new computed
// values, and how, the moment the cascade changes them. Animations is the
// finished animation list, the keyframes blocks the element plays and how.
// Both are assembled from the parallel longhand lists by
// finishTransitions and finishAnimations at the end of the cascade; the
// longhand lists are intermediate and never meant for the template.
type motionStyle struct {
	Transitions []Transition
	Animations  []Animation

	TransitionProps []string
	TransitionDurs  []Time
	TransitionTims  []Timing
	TransitionDels  []Time
	TransitionNone  bool

	AnimationNames []string
	AnimationDurs  []Time
	AnimationTims  []Timing
	AnimationDels  []Time
	AnimationIters []float64 // math.Inf(+1) means infinite
	AnimationDirs  []uint8
	AnimationFills []uint8
}
