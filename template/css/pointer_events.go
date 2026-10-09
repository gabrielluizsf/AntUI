package css

// PointerEvents says whether a box takes part in hit testing.
const (
	PointerEventsAuto uint8 = iota
	PointerEventsNone
)

func parsePointerEvents(raw string) (uint8, bool) {
	switch raw {
	case "auto":
		return PointerEventsAuto, true
	case "none":
		return PointerEventsNone, true
	}
	return 0, false
}
