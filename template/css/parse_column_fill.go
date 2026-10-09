package css

import (
	"strings"
)

// parseColumnFill reads how a declared height is spent across the columns.
func parseColumnFill(raw string) (uint8, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "balance":
		return ColumnFillBalance, true
	case "auto":
		return ColumnFillAuto, true
	}
	return 0, false
}
