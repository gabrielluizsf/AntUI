package svg

// filterDef is a `<filter>` element: the primitives it holds, in the order
// they were written, which is the order they run in.
type filterDef struct {
	prims []filterPrim
}
