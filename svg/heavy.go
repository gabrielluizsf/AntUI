package svg

// heavy is a weight the canvas draws heavier than the face it has: a face is
// cut in the weights it is cut in, and the ones it does not have are faked by
// drawing the writing a second time a pixel on.
func heavy(weight uint16) bool {
	if weight == 0 {
		weight = FontWeightNormal
	}
	return weight >= 600
}
