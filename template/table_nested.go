package template

// tableNested is one table drawn inside another container. It carries the batch
// that lines its rows up — collected by the silent pass and re-solved at the
// final box — and the border box the table asks its parent for.
type tableNested struct {
	batch  *tableBatch
	nw, nh int
}
