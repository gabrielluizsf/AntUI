package css

import (
	"math"
	"strconv"
	"strings"
)

// Keyframe is one stop of an @keyframes block: the progress at which it
// fires — 0 for from, 1 for to, otherwise the percentage — and the
// declarations it carries. The declarations apply to a base style the same
// way a rule's do, so a keyframe can mention any property the engine knows.
type Keyframe struct {
	Offset float64 // 0..1
	Decls  []Declaration
}

// Keyframes is a parsed @keyframes block, in the order its author wrote the
// frames. Style.Animations name it, and [Animate] walks its frames with the
// progress of an animation.
type Keyframes struct {
	Name   string // the name @keyframes got, which animation-name spells
	Frames []Keyframe
}

// Keyframes returns the @keyframes block with the given name, or nil when
// the sheet never defined it.
func (sh *Sheet) Keyframes(name string) *Keyframes {
	if sh.keyframes == nil {
		return nil
	}
	return sh.keyframes[name]
}

// addKeyframes stores one parsed block, a later definition of the same name
// replacing the earlier one the way CSS reads the last @keyframes it met.
func (sh *Sheet) addKeyframes(kf *Keyframes) {
	if sh.keyframes == nil {
		sh.keyframes = map[string]*Keyframes{}
	}
	sh.keyframes[kf.Name] = kf
}

// parseKeyframes reads the body of an @keyframes block: a series of frames
// introduced by from, to or a percentage followed by a declaration block.
// A frame with an unrecognised main is consumed and skipped, and a malformed
// body ends the rule the way any structural damage does.
func (p *parser) parseKeyframes(sh *Sheet, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		// A nameless @keyframes is nothing a rule can name; the block is
		// slurped and dropped.
		return p.skipBlock()
	}
	kf := &Keyframes{Name: name}
	for {
		p.skipSpace()
		if p.eof() {
			return fmtErrf("unterminated @keyframes %q, missing '}'", name)
		}
		if p.peek() == '}' {
			p.next()
			sh.addKeyframes(kf)
			return nil
		}
		header, semi, err := p.readHeader(true)
		if err != nil {
			return err
		}
		off, ok := parseKeyframeOffset(header)
		if !ok {
			if semi {
				if p.peek() == ';' {
					p.next()
				}
			} else if p.peek() == '{' {
				if err := p.skipBlock(); err != nil {
					return err
				}
			}
			continue
		}
		if p.peek() == '{' {
			p.next()
		}
		kf.Frames = append(kf.Frames, Keyframe{Offset: off, Decls: p.readFrameDecls()})
	}
}

// readFrameDecls reads one keyframe's declaration block, returning the
// declarations it held. A missing '}' ends the block and the parse, exactly
// as a rule that never closes would.
func (p *parser) readFrameDecls() []Declaration {
	var decls []Declaration
	for {
		p.skipSpace()
		if p.eof() {
			return decls
		}
		if p.peek() == '}' {
			p.next()
			return decls
		}
		text, err := p.readDecl()
		if err != nil {
			return decls
		}
		if p.peek() == ';' {
			p.next()
		}
		if d, ok := parseDeclaration(text); ok {
			decls = append(decls, d)
		}
	}
}

// parseKeyframeOffset reads a frame's main: from, to, or a percentage such
// as 32.5%. Percentages clamp into [0,1]; junk reports failure.
func parseKeyframeOffset(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "from":
		return 0, true
	case "to":
		return 1, true
	}
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, false
		}
		return clamp01(v / 100), true
	}
	return 0, false
}

// frameStyle computes the style one keyframe paints: its declarations folded
// into the base's custom-property map, so var(--x) reads the element's own
// values. The frame style carries only what its own Set map saw.
func frameStyle(frame Keyframe, base Style, ctx Units) Style {
	var st Style
	st.Custom = base.Custom
	st.inherit = map[string]bool{}
	set := map[string]bool{}
	for _, d := range frame.Decls {
		applyDecl(&st, set, base.Custom, base.Custom, d, &ctx)
	}
	st.Set = set
	resolveCurrentColors(&st)
	return st
}

// frameSpan finds the two frames straddling a progress t in [0,1], with the
// local progress between them. t before the first frame pins to the first,
// after the last pins to the last.
func frameSpan(frames []Keyframe, t float64) (lo, hi Keyframe, local float64, ok bool) {
	if len(frames) == 0 {
		return Keyframe{}, Keyframe{}, 0, false
	}
	if t <= frames[0].Offset || len(frames) == 1 {
		return frames[0], frames[0], 0, true
	}
	last := frames[len(frames)-1]
	if t >= last.Offset {
		return last, last, 1, true
	}
	for i := 0; i < len(frames)-1; i++ {
		a, b := frames[i], frames[i+1]
		if t >= a.Offset && t <= b.Offset {
			local := 0.0
			if span := b.Offset - a.Offset; span > 0 {
				local = (t - a.Offset) / span
			}
			return a, b, local, true
		}
	}
	return frames[0], frames[0], 0, true
}

// Animate applies a @keyframes block at the progress t (0..1) over the base
// style and returns the animated style. The result owns its Set map, so the
// cached base style is never touched; properties the frames wrote overlay
// the base, and between two frames every animatable property the frames
// mentioned interpolates, with the base style's value standing in for a
// frame that omitted it. Properties the engine does not interpolate jump at
// the 0.5 crossing.
func Animate(base Style, kf *Keyframes, t float64, ctx Units) Style {
	if kf == nil || len(kf.Frames) == 0 {
		return base
	}
	out := CopyStyle(base)
	lo, hi, local, ok := frameSpan(kf.Frames, t)
	if !ok {
		return out
	}
	if lo.Offset == hi.Offset {
		fa := frameStyle(lo, base, ctx)
		overlayProps(&out, fa, fa.Set)
		return out
	}
	fa := frameStyle(lo, base, ctx)
	fb := frameStyle(hi, base, ctx)
	props := map[string]bool{}
	for _, prop := range animatableProps {
		if fa.Set[prop] || fb.Set[prop] {
			props[prop] = true
		}
	}
	BlendKeyframes(&out, base, fa, fb, props, local, ctx)
	return out
}
