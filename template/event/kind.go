package event

// Kind is what happened to a component this frame.
type Kind int

// The things a component can report.
const (
	None   Kind = iota
	Click       // a button was pressed all the way down and released on it
	Toggle      // a checkbox flipped
	Pick        // a radio was picked, or a dropdown option was chosen
	Change      // a slider moved
	Type        // text was typed into or deleted from a field
	Focus       // a text field got the keyboard
	Open        // a dropdown opened
	Close       // a dropdown closed
)
