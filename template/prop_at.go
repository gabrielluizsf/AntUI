package template

// propAt is the layer-list cycling: a single declared value applies to every
// layer, the next value to the next layer, wrapping around.
func propAt[T any](list []T, i int, def T) T {
	if len(list) == 0 {
		return def
	}
	return list[i%len(list)]
}
