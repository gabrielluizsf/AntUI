package template

// widgetBox is a widget's place and its identity: the box the focus ring hangs
// on, and the box the pointer counts as being over it.
//
// The two are usually the same, and where they are not — a checkbox whose label
// is part of what can be clicked but not part of what the keyboard focuses —
// each is given here rather than assumed. kind is what the widget calls itself
// to the window, and label goes into the identity so that two identical
// controls on one page stay two widgets.
type widgetBox struct {
	kind, label string
	x, y        int
	idW, idH    int // the box the identity is made from
	hitW, hitH  int // the box the pointer counts as being over it
}
