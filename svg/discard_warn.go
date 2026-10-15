package svg

// discardWarn is the warn a style is read with where the reading is only for the
// style itself: an attribute that cannot be read is said where the element is
// drawn, and saying it a second time here would only repeat it.
func discardWarn(string, ...any) {}
