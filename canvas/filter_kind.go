package canvas

// FilterKind names one CSS filter function the canvas can apply to a region.
// The values are deliberately independent of any stylesheet; a template maps
// its own parsed filters onto them.
type FilterKind uint8
