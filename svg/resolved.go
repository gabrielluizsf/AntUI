package svg

// resolved is the gradient as it will be painted: every attribute it does not
// say taken from the gradient its href points at, which may itself point at
// another. The chain is walked once with the whole drawing in hand, which is why
// an href cannot be followed by [parseGradient] on its own.
//
// A gradient that points at itself, or in a circle, stops where it would start
// repeating: everything such a chain has to say is already in it by then, and
// following it further would loop.
//
// A radial gradient that never said where its focus is has it at the centre,
// which is only known once the chain above has settled the centre itself.
func (g *gradient) resolved(all map[string]*gradient) *gradient {
	out := gradientDefaults(g.radial)
	out.viewport = g.viewport
	var chain []*gradient
	for at, seen := g, map[*gradient]bool{}; !seen[at]; {
		seen[at] = true
		chain = append(chain, at)
		if at.href == "" {
			break
		}
		next, ok := all[at.href]
		if !ok {
			break
		}
		at = next
	}
	// The chain runs nearest first, and the first gradient to have said something
	// is the one that counts: an href fills in what was left out and nothing more.
	for _, at := range chain {
		out.inherit(at)
	}
	if g.radial {
		if !out.set["fx"] {
			out.fx = out.cx
		}
		if !out.set["fy"] {
			out.fy = out.cy
		}
	}
	out.resolvePercentages()
	return out
}
