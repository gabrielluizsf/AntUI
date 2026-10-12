//go:build android

package sqlite

// The types a column can hold, as Cursor numbers them.
const (
	typeNull = 0
	typeInt  = 1
	typeReal = 2
	typeText = 3
	typeBlob = 4
)
