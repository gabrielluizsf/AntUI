//go:build android

package gallery

// The columns read. They are the same names on every kind of media, because
// they are MediaStore.MediaColumns and not the per-kind ones.
var listColumns = []string{
	"_id", "_display_name", "_size", "mime_type", "date_added", "duration",
}
