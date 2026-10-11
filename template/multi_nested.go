package template

// multiNested is one multi-column container drawn inside another container. It
// carries the batch that lays the container's children out — collected by the
// silent pass and re-solved at the final box — and the border box the container
// asks its parent for.
type multiNested struct {
	batch  *multiBatch
	nw, nh int
}
