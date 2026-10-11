package canvas

// Fit is the largest width and height with the source's proportions that
// fits inside the box given.
func Fit(srcW, srcH, boxW, boxH int) (width, height int) {
	if srcW <= 0 || srcH <= 0 || boxW <= 0 || boxH <= 0 {
		return 0, 0
	}
	if srcW*boxH < boxW*srcH {
		return srcW * boxH / srcH, boxH
	}
	return boxW, srcH * boxW / srcW
}
