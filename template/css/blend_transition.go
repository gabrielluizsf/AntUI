package css

// blendTransition is the style a transition draws at a moment in time: the
// target [Style] is blended in from the entry [Style] for every transition
// it declared. elapsed is milliseconds since the transition began. While a
// transition waits out its delay the entry's value shows; past its end the
// target shows, so a finished transition is the target itself.
type BlendTransitionParams struct {
	From    Style
	To      Style
	Elapsed float64
	Units   Units
}

func BlendTransition(p BlendTransitionParams) Style {
	from, to, elapsed, ctx := p.From, p.To, p.Elapsed, p.Units
	out := CopyStyle(to)
	for _, tx := range to.Transitions {
		start := tx.Delay.MS()
		dur := tx.Duration.MS()
		if dur <= 0 {
			continue
		}
		if tx.Prop == "all" {
			for _, prop := range animatableProps {
				blendPropTimed(&out, from, prop, start, dur, elapsed, tx.Timing, ctx)
			}
			continue
		}
		if tx.Prop == "none" {
			continue
		}
		blendPropTimed(&out, from, tx.Prop, start, dur, elapsed, tx.Timing, ctx)
	}
	return out
}

func blendTransition(from, to Style, elapsed float64, ctx Units) Style {
	return BlendTransition(BlendTransitionParams{From: from, To: to, Elapsed: elapsed, Units: ctx})
}
