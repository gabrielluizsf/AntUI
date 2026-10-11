package template

// gridNested is one grid container drawn inside another container. It carries
// the batch that lays the container's children out — collected by the silent
// pass and re-solved at the final box — and the border box the container asks
// its parent for.
type gridNested struct {
	batch      *gridBatch
	nw, nh     int
	minW, minH int
}
