//go:build android

package app

// OK reports whether the other activity finished rather than being backed
// out of.
func (r ActivityResult) OK() bool { return r.Code == ResultOK }
