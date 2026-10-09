package css

// markUnsupported records the first reason a selector cannot be evaluated,
// keeping the most specific one.
func markUnsupported(prev, next string) string {
	if prev != "" {
		return prev
	}
	return next
}
