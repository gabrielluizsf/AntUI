package font

// hMetrics reads the horizontal half of a metrics record: the same five
// bytes whether the record is the small one or the hori half of the big one.
func hMetrics(d []byte) sbitMetrics {
	if len(d) < 5 {
		return sbitMetrics{}
	}
	return sbitMetrics{
		Height:   int(d[0]),
		Width:    int(d[1]),
		BearingX: int(int8(d[2])),
		BearingY: int(int8(d[3])),
		Advance:  int(d[4]),
	}
}
