package css

// parseAlignSelf turns the align-self keywords into an Align* constant,
// auto included — the only keyword align-items rejects.
func parseAlignSelf(raw string) (uint8, bool) {
	if raw == "auto" {
		return AlignAuto, true
	}
	return parseAlignItems(raw)
}
