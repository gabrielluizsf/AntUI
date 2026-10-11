package antui

// widgetID hashes a widget's identity with FNV-1a.
func widgetID(kind string, x, y, w, h int, label string) uint32 {
	const (
		offset = 2166136261
		prime  = 16777619
	)
	hash := uint32(offset)
	mix := func(b byte) {
		hash ^= uint32(b)
		hash *= prime
	}
	for _, v := range [4]int{x, y, w, h} {
		for shift := range 4 {
			mix(byte(uint32(v) >> (shift * 8)))
		}
	}
	for i := range len(kind) {
		mix(kind[i])
	}
	for i := range len(label) {
		mix(label[i])
	}
	if hash == 0 {
		return 1
	}
	return hash
}
