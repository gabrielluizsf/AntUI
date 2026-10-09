package css

// Keyframes is a parsed @keyframes block, in the order its author wrote the
// frames. Style.Animations name it, and [Animate] walks its frames with the
// progress of an animation.
type Keyframes struct {
	Name   string // the name @keyframes got, which animation-name spells
	Frames []Keyframe
}
