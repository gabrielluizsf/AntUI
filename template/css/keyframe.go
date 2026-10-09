package css

// Keyframe is one stop of an @keyframes block: the progress at which it
// fires — 0 for from, 1 for to, otherwise the percentage — and the
// declarations it carries. The declarations apply to a base style the same
// way a rule's do, so a keyframe can mention any property the engine knows.
type Keyframe struct {
	Offset float64 // 0..1
	Decls  []Declaration
}
