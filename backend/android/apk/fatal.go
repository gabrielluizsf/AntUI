package apk

// Fatal reports whether any of these would stop an upload.
func Fatal(ps []Problem) bool {
	for _, p := range ps {
		if p.Fatal {
			return true
		}
	}
	return false
}
