package antui

// Rate is how many frames a second the device is really running at.
func (a *Audio) Rate() int {
	if a == nil {
		return 0
	}
	return a.rate
}
