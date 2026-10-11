package template

// coverSize is the smallest scale with the picture's proportions that still
// covers the box: the height scaled to the box height when that is the tight
// fit, the width otherwise.
func coverSize(srcW, srcH, boxW, boxH int) (int, int) {
	if srcW <= 0 || srcH <= 0 || boxW <= 0 || boxH <= 0 {
		return 0, 0
	}
	if srcW*boxH > boxW*srcH {
		return srcW * boxH / srcH, boxH
	}
	return boxW, srcH * boxW / srcW
}
