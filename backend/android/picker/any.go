//go:build android

package picker

// Any asks for a file of a given type — "application/pdf", "text/*", or
// "*/*" for anything at all.
func Any(mime string) (string, error) { return one(mime) }
