//go:build android

package display

// Refresh is how many frames a second the app's screen shows. It is what
// antui asks for when it wants to know the frame budget.
func Refresh() (float32, error) {
	d, err := Default()
	return d.Refresh, err
}
