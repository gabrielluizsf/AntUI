package audio

// normalize lifts a quiet recording up to a full swing, so a bank's sounds
// all come out at the same presence whatever the source level was.
func normalize(samples []float32) {
	var peak float32
	for _, s := range samples {
		if a := abs(s); a > peak {
			peak = a
		}
	}
	if peak == 0 || peak >= 0.99 {
		return
	}
	gain := 0.92 / peak
	for i := range samples {
		samples[i] *= gain
	}
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
