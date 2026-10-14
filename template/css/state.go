package css

type State uint8

const (
	StateNone  State = 0
	StateHover State = 1 << iota
	StateFocus
	StateActive
	StateChecked
)

// StateFlags carries the boolean interaction flags for building a State.
type StateFlags struct {
	Hovered bool
	Focused bool
	Active  bool
	Checked bool
}

// StateWith builds a State for the interaction layer's flags.
func StateWith(flags StateFlags) State {
	var s State
	if flags.Hovered {
		s |= StateHover
	}
	if flags.Focused {
		s |= StateFocus
	}
	if flags.Active {
		s |= StateActive
	}
	if flags.Checked {
		s |= StateChecked
	}
	return s
}
