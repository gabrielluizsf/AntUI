package apk

// String is the problem as one line, or two when it says what to do
// about it.
func (p Problem) String() string {
	mark := "warning"
	if p.Fatal {
		mark = "refused"
	}
	if p.Fix == "" {
		return mark + ": " + p.Text
	}
	return mark + ": " + p.Text + "\n         " + p.Fix
}
