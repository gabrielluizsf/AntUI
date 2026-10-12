package sdk

// pick chooses the newest of what was found, preferring a release over any
// preview however new the preview is — an rc of the build tools has shipped
// broken often enough that it is not worth taking by accident.
func pick[T any](found []T, version func(T) Version) (T, bool) {
	var best T
	var bestV Version
	ok := false
	for _, f := range found {
		v := version(f)
		switch {
		case !ok:
		case bestV.Release() && !v.Release():
			continue
		case !bestV.Release() && v.Release():
		case v.Compare(bestV) <= 0:
			continue
		}
		best, bestV, ok = f, v, true
	}
	return best, ok
}
