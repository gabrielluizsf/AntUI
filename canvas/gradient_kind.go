package canvas

// GradientKind names one CSS gradient the canvas can fill a rectangle with.
// The values are deliberately independent of any stylesheet; a template maps
// its own parsed gradients onto them.
type GradientKind uint8
