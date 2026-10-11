package template

// A Screen is one view in a window: it draws itself on a Context and decides
// on its own when to leave. A screen knows nothing about the screens around
// it — the Context holds them, as a stack with the showing screen on top —
// so screens serve as the pointers a game navigates along:
//
//	home → shop → cart, and back the same way.
type Screen interface {
	Draw(c *Context)
}
