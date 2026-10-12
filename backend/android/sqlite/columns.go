//go:build android

package sqlite

// Columns is what the query answered with.
func (r *Rows) Columns() []string { return r.columns }
