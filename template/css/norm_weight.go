package css

func normWeight(w uint16) uint16 {
	if w == 0 {
		return FontWeightNormal
	}
	return min(max((w+50)/100*100, 100), maxWeight)
}
