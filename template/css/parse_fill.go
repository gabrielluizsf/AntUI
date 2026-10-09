package css

// parseFill reads one fill-mode keyword.
func parseFill(raw string) (uint8, bool) {
	switch raw {
	case "none":
		return FillNone, true
	case "backwards":
		return FillBackwards, true
	case "forwards":
		return FillForwards, true
	case "both":
		return FillBoth, true
	}
	return 0, false
}
