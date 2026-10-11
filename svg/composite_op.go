package svg

// compositeOp is feComposite's operator; the zero is "over", which is what the
// primitive means when the operator is left off.
type compositeOp uint8
