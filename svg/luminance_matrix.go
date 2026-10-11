package svg

func luminanceMatrix() [20]float64 {
	return [20]float64{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0.2125, 0.7154, 0.0721, 0, 0,
	}
}
