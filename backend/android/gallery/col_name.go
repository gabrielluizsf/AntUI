//go:build android

package gallery

// The MediaStore columns used here, by the names the platform gives them.
const (
	colName    = "_display_name"
	colMime    = "mime_type"
	colPath    = "relative_path" // API 29 and up
	colPending = "is_pending"    // API 29 and up
)
