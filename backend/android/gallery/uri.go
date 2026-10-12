//go:build android

package gallery

// URI is the content:// address of what is being written. It is what to hand
// another app, and it is not a path: on a modern device there may be no path
// at all.
func (w *Writer) URI() string { return w.uriText }
