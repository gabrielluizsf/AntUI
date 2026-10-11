package font

// sbitMetrics is a bitmap's size and placement in the pixels of its strike.
// The horizontal bearings of a big metrics record are what a run of text
// uses, and they are its first five bytes — which is all a small record
// holds.
type sbitMetrics struct {
	Height             int
	Width              int
	BearingX, BearingY int
	Advance            int
}
