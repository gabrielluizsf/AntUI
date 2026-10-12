//go:build android

package gallery

// Query is how to narrow a listing.
type Query struct {
	// Album keeps only what is in one folder — the same name [Create] took.
	// Empty means everything.
	Album string
	// Limit is how many at most, newest first. Zero means everything, which
	// on a real phone can be tens of thousands.
	Limit int
	// Mine keeps only what this app put there. It is the one listing that
	// needs no permission at all: an app can always see its own.
	Mine bool
}
