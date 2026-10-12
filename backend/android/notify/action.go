//go:build android

package notify

// Action is a button on a notification.
//
// It carries a **name**, not a function, and that is not an inconvenience —
// it is the only thing that works. A notification outlives the app that
// posted it: the user may tap the button an hour later, on a device that
// killed the process long ago, and Android starts the app afresh to deliver
// it. A Go closure cannot survive that. A name can, and [OnAction] is where
// the app says what each name means, at startup, every time it starts.
type Action struct {
	// Title is the label on the button.
	Title string
	// Name is what [OnAction] matched it against. It never appears on screen.
	Name string
	// Icon is a drawable resource id, and may be zero: from Android 7 the
	// icon on an action button is not drawn at all, and it is only still
	// required because the API predates that.
	Icon int32
	// Reply turns the button into a text field the user can type into, which
	// is how a message is answered from the shade.
	Reply *Reply
}
