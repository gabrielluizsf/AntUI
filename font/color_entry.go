package font

// colorEntry reads one CBDT image: the metrics it carries, if any, and the
// PNG behind its dataLen. The three formats differ only in which metrics sit
// in front of the length.
func colorEntry(b colorBitmap, d []byte) (sbitMetrics, []byte, bool) {
	var m sbitMetrics
	var at int
	switch b.imageFormat {
	case 17: // small metrics, then the length, then the PNG
		m, at = hMetrics(d), 9
	case 18: // big metrics, then the length, then the PNG
		m, at = hMetrics(d), 12
	case 19: // the length and the PNG; the metrics are the index's
		m, at = b.metrics, 4
		if !b.inIndex {
			return m, nil, false
		}
	default:
		return m, nil, false
	}
	if at > len(d) {
		return m, nil, false
	}
	n := int(be32(d, at-4))
	if n < 0 || at+n > len(d) {
		return m, nil, false
	}
	return m, d[at : at+n], true
}
