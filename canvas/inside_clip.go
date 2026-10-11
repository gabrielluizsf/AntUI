package canvas

func (cv *Canvas) insideClip(x, y int) bool {
	return x >= cv.Clip.X && y >= cv.Clip.Y &&
		x < cv.Clip.X+cv.Clip.Width && y < cv.Clip.Y+cv.Clip.Height
}
