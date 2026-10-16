package svg

// Warnings is everything the drawing asked for that this package could not do,
// read after it is painted. A drawing that came out as it should has none.
//
// A picture read from a server is said to have failed from the goroutine that
// fetched it, so a drawing whose pictures come over the network can gain one
// after the call that read it returned — which is the moment the photo was
// expected and did not come. The list is copied out under the lock that guards
// it, so the answer is the moment it was asked for rather than a slice another
// goroutine may be appending to.
func (img *Image) Warnings() Warnings {
	img.warnMu.Lock()
	defer img.warnMu.Unlock()
	return append(Warnings(nil), img.warnings...)
}
