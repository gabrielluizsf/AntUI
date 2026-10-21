package antui

// RadioParams groups a radio option's label, bound value and option number.
// A value of NoOption means no option in the group is selected.
type RadioParams struct {
	Label  string
	Value  *int
	Option int
}

// NoOption is the value a radio group holds while nothing in it is picked.
// Clicking the picked option clears the group back to it.
const NoOption = -1
