//go:build android

package ndk

import "fmt"

// sprintf is fmt.Sprintf behind a name, so that the logging calls above read
// without an import that looks like it might be doing something else.
func sprintf(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}
