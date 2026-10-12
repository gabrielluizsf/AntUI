package apk

func plural(n int) string {
	if n == 1 {
		return "y is"
	}
	return "ies are"
}
