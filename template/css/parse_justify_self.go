package css

// parseJustifySelf reads the justify-self keywords for one cell.
func parseJustifySelf(raw string) (uint8, bool) {
	if raw == "auto" {
		return AlignAuto, true
	}
	return parseJustifyItems(raw)
}

// gridFunction splits a track function call — "minmax(0, 1fr)", "repeat(2, …)",
// "fit-content(200px)" — into its lowercased name and its arguments. A function
// name is case-insensitive, so it is folded for the caller to match.
