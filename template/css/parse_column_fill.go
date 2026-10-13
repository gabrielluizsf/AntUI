package css

// parseColumnFill reads how a declared height is spent across the columns.
func parseColumnFill(raw string) (uint8, bool) {
	switch keyword(raw) {
	case "balance":
		return ColumnFillBalance, true
	case "auto":
		return ColumnFillAuto, true
	}
	return 0, false
}
