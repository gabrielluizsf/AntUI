//go:build android

package sqlite

// Err is what stopped the walk, or nil.
func (r *Rows) Err() error { return r.err }
