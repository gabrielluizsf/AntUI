package canvas

// blurBuffers is the pair of float buffers a blur moves its picture between,
// sized to n floats and kept on the canvas so the next blur of the same
// canvas — next frame, most likely — finds them already there.
func (cv *Canvas) blurBuffers(n int) (front, back []float64) {
	if cap(cv.blurFront) < n {
		cv.blurFront = make([]float64, n)
		cv.blurBack = make([]float64, n)
	}
	return cv.blurFront[:n], cv.blurBack[:n]
}
