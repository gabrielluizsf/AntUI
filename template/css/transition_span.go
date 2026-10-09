package css

// TransitionSpan is the moment the last transition in a style reaches its
// target: the longest duration-plus-delay. A style whose transitions all run
// instantly reports zero, which a caller reads as "snap, do not animate".
func TransitionSpan(st Style) float64 {
	span := 0.0
	for _, tx := range st.Transitions {
		if end := tx.Delay.MS() + tx.Duration.MS(); end > span {
			span = end
		}
	}
	return span
}
