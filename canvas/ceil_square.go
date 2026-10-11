package canvas

// ceilSquare is ⌈v²/65536⌉, the smallest distance whose fixed-point root is at
// least v — the threshold a squared distance is measured against to decide
// which side of an edge a pixel is on without taking its root.
func ceilSquare(v uint64) uint64 { return (v*v + 0xFFFF) >> 16 }
