package css

import (
	"math"
	"strconv"
	"strings"
)

// The direction an animation travels along: whether each iteration runs
// forward, backward, or alternates between them. The zero is normal.
const (
	AnimNormal uint8 = iota
	AnimReverse
	AnimAlternate
	AnimAlternateReverse
)

// The fill-mode of an animation: whether the frames bleed over the delay and
// the end, so a button that animates once keeps its final colour. The zero
// is none.
const (
	FillNone uint8 = iota
	FillBackwards
	FillForwards
	FillBoth
)

// Animation is one entry of the animation shorthand: the @keyframes block it
// plays (by name), its duration, easing and delay, how many times it runs
// (Infinite) and which direction and fill-mode govern it. Iterations without
// Infinite is the written number, defaulting to 1.
type Animation struct {
	Name       string
	Duration   Time
	Timing     Timing
	Delay      Time
	Iterations float64
	Infinite   bool
	Direction  uint8 // one of the Anim* constants
	Fill       uint8 // one of the Fill* constants
}

// parseAnimations reads the animation shorthand: a comma-separated list of
// animations, each made of a keyframe name, two times, an easing, an
// iteration count, and the direction and fill keywords in any order. An entry
// with no name — an empty list, a stray keyword — invalidates the whole
// property, exactly as a browser drops an unparsable declaration.
func parseAnimations(raw string) ([]Animation, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []Animation
	for _, part := range splitFields(raw, ',') {
		a, ok := parseAnimation(strings.TrimSpace(part))
		if !ok {
			return nil, false
		}
		out = append(out, a)
	}
	return out, len(out) > 0
}

// parseAnimation reads one comma-free animation. Names are kept as the
// author wrote them; the keywords are case-insensitive.
func parseAnimation(part string) (Animation, bool) {
	part = strings.TrimSpace(part)
	if part == "" {
		return Animation{}, false
	}
	a := Animation{Iterations: 1, Direction: AnimNormal, Fill: FillNone}
	var sawDur, sawDelay bool
	for _, tk := range splitTokens(part) {
		t := strings.ToLower(tk)
		switch t {
		case "infinite":
			a.Infinite = true
		case "reverse":
			a.Direction = AnimReverse
		case "alternate":
			a.Direction = AnimAlternate
		case "alternate-reverse":
			a.Direction = AnimAlternateReverse
		case "backwards":
			a.Fill = FillBackwards
		case "forwards":
			a.Fill = FillForwards
		case "both":
			a.Fill = FillBoth
		case "none":
			// fill-mode: none, the default.
		default:
			if tim, ok := parseTiming(t); ok {
				a.Timing = tim
				continue
			}
			if tm, ok := ParseTime(t); ok {
				if !sawDur {
					a.Duration = tm
					sawDur = true
				} else if !sawDelay {
					a.Delay = tm
					sawDelay = true
				} else {
					return Animation{}, false
				}
				continue
			}
			if v, err := strconv.ParseFloat(t, 64); err == nil {
				a.Iterations = v
				continue
			}
			if strings.ContainsRune(t, '(') {
				return Animation{}, false
			}
			if a.Name != "" {
				return Animation{}, false
			}
			a.Name = trimQuotes(tk)
		}
	}
	return a, a.Name != ""
}

// trimQuotes peels the quotes CSS allows around a custom-ident or string.
func trimQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '"' && s[len(s)-1] == '"' || s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

// animationParts splits a parsed shorthand into the parallel lists the
// longhands also fill, so the cascade merges them in one finisher. An
// infinite count lands in the list as positive infinity.
func animationParts(as []Animation) (names []string, durs []Time, tims []Timing, dels []Time, iters []float64, dirs, fills []uint8) {
	for _, a := range as {
		names = append(names, a.Name)
		durs = append(durs, a.Duration)
		tims = append(tims, a.Timing)
		dels = append(dels, a.Delay)
		if a.Infinite {
			iters = append(iters, math.Inf(1))
		} else {
			iters = append(iters, a.Iterations)
		}
		dirs = append(dirs, a.Direction)
		fills = append(fills, a.Fill)
	}
	return names, durs, tims, dels, iters, dirs, fills
}

// finishAnimations merges the animation longhands (and the shorthand's own
// parallel lists) into the style's final Animations slice. The name list
// sets the length, and — like CSS — every shorter list cycles to cover the
// whole row, or leaves its CSS initial where the longhand never appeared.
func finishAnimations(st *Style) {
	if len(st.AnimationNames) == 0 {
		st.Animations = nil
		return
	}
	for i, name := range st.AnimationNames {
		a := Animation{Name: name, Iterations: 1, Direction: AnimNormal, Fill: FillNone}
		if len(st.AnimationDurs) > 0 {
			a.Duration = st.AnimationDurs[i%len(st.AnimationDurs)]
		}
		if len(st.AnimationTims) > 0 {
			a.Timing = st.AnimationTims[i%len(st.AnimationTims)]
		}
		if len(st.AnimationDels) > 0 {
			a.Delay = st.AnimationDels[i%len(st.AnimationDels)]
		}
		if len(st.AnimationIters) > 0 {
			if v := st.AnimationIters[i%len(st.AnimationIters)]; math.IsInf(v, 1) {
				a.Infinite = true
			} else {
				a.Iterations = v
			}
		}
		if len(st.AnimationDirs) > 0 {
			a.Direction = st.AnimationDirs[i%len(st.AnimationDirs)]
		}
		if len(st.AnimationFills) > 0 {
			a.Fill = st.AnimationFills[i%len(st.AnimationFills)]
		}
		st.Animations = append(st.Animations, a)
	}
}

// parseNameList reads an animation-name value: a comma-separated series of
// keyframe names.
func parseNameList(raw string) ([]string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []string
	for _, p := range splitFields(raw, ',') {
		p = trimQuotes(strings.TrimSpace(p))
		if p == "" {
			return nil, false
		}
		out = append(out, p)
	}
	return out, true
}

// parseIterations reads one iteration count: a number or the infinite
// keyword, which lands as positive infinity.
func parseIterations(raw string) (float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "infinite" {
		return math.Inf(1), true
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

// parseIterationList reads a comma-separated list of iteration counts.
func parseIterationList(raw string) ([]float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []float64
	for _, p := range splitFields(raw, ',') {
		v, ok := parseIterations(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}

// parseDirection reads one direction keyword.
func parseDirection(raw string) (uint8, bool) {
	switch raw {
	case "normal":
		return AnimNormal, true
	case "reverse":
		return AnimReverse, true
	case "alternate":
		return AnimAlternate, true
	case "alternate-reverse":
		return AnimAlternateReverse, true
	}
	return 0, false
}

// parseDirectionList reads a comma-separated list of direction keywords.
func parseDirectionList(raw string) ([]uint8, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []uint8
	for _, p := range splitFields(raw, ',') {
		v, ok := parseDirection(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}

// parseFill reads one fill-mode keyword.
func parseFill(raw string) (uint8, bool) {
	switch raw {
	case "none":
		return FillNone, true
	case "backwards":
		return FillBackwards, true
	case "forwards":
		return FillForwards, true
	case "both":
		return FillBoth, true
	}
	return 0, false
}

// parseFillList reads a comma-separated list of fill-mode keywords.
func parseFillList(raw string) ([]uint8, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	var out []uint8
	for _, p := range splitFields(raw, ',') {
		v, ok := parseFill(strings.TrimSpace(p))
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}
