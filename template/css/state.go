package css

type State uint8

const (
	StateNone  State = 0
	StateHover State = 1 << iota
	StateFocus
	StateActive
	StateChecked
)

// StateWith builds a State for the interaction layer's flags.
func StateWith(hovered, focused, active, checked bool) State {
	var s State
	if hovered {
		s |= StateHover
	}
	if focused {
		s |= StateFocus
	}
	if active {
		s |= StateActive
	}
	if checked {
		s |= StateChecked
	}
	return s
}
