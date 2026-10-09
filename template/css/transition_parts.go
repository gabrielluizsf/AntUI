package css

// transitionParts splits a parsed shorthand into the four parallel lists the
// longhands also fill, so the cascade can merge them all in one finisher.
func transitionParts(ts []Transition) (props []string, durs []Time, tims []Timing, dels []Time) {
	for _, tr := range ts {
		props = append(props, tr.Prop)
		durs = append(durs, tr.Duration)
		tims = append(tims, tr.Timing)
		dels = append(dels, tr.Delay)
	}
	return props, durs, tims, dels
}
