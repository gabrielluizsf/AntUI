package svg

// namedColors is the set of colour names a drawing may use in place of a hex
// value, which is the whole of the CSS named colours plus `transparent`. The
// keys are lowercased, because a name is matched without regard to case.
var namedColors = buildNamedColors()
