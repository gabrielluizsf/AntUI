//go:build android

package dialog

// Which button was pressed. The numbers are the platform's, and they are
// negative — DialogInterface counts its buttons downwards from -1.
type Button int
