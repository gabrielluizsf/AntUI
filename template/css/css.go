package css

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// State is a widget's interaction state, the value pseudo-classes read.
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

// Selector is one simple selector that must match for a rule to apply.
// `.a.b:hover` carries the classes "a" and "b" and the Hover state. Descendant
// and child combinators are parsed but can never match: AntUI widgets are
// flat, so the engine keeps the rule with a warning and the selector never
// fires. The same goes for parts of the CSS3 selector language the flat
// model cannot evaluate — attributes, ids, pseudo-elements and structural
// pseudo-classes — all recorded in [Selector.Unsupported].
type Selector struct {
	Tag     string   // element name; "" means any element with the classes
	Classes []string // every class must be present
	State   State    // every bit must be present
	All     bool     // the universal selector *

	// Unsupported names the selector feature the flat widget model cannot
	// evaluate (a combinator, an attribute, an id, a pseudo-element, a
	// structural pseudo-class). Non-empty means the rule parsed fine but must
	// never match — the engine drops it with a warning rather than an error.
	Unsupported string
}

func (s Selector) match(tag string, classes []string, state State) bool {
	if s.Unsupported != "" {
		return false
	}
	tag = strings.ToLower(tag)
	if s.All {
		if s.Tag != "" && s.Tag != tag {
			return false
		}
	} else if s.Tag != "" && s.Tag != tag {
		return false
	}
	for _, c := range s.Classes {
		found := false
		for _, have := range classes {
			if have == strings.ToLower(c) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return s.State&state == s.State
}

// ClassList splits a space-separated set of class names.
func ClassList(s string) []string {
	return strings.Fields(s)
}

// specificity returns the selector's [a, b, c] triple: classes and
// pseudo-classes, elements, universals. A higher triple wins the cascade.
func (s Selector) specificity() (a, b, c int) {
	b += len(s.Classes)
	for st := State(1); st <= StateChecked; st <<= 1 {
		if s.State&st != 0 {
			b++
		}
	}
	if !s.All && s.Tag != "" {
		c = 1
	}
	return a, b, c
}

// Media is a window-width window a rule lives inside. A stylesheet may list
// several media queries — comma- and or-separated — and the rule applies while
// any of them holds. Only the width features (min-width/max-width) are
// evaluated; every other condition (orientation, resolution, a media type such
// as print) is parsed, warned about, and treated as satisfied so the rule is
// never dropped on a canvas. Negation flips a query's width window.
type Media struct {
	Queries []MediaQuery
}

// MediaQuery is one alternative of an @media list: a width window, optional
// negation. A query with no window matches any width; that is also how
// conditions the engine cannot evaluate fall out.
type MediaQuery struct {
	MinWidth int
	MaxWidth int
	HasMin   bool
	HasMax   bool
	Negated  bool
}

// matches tests one alternative against the window width.
func (q MediaQuery) matches(width int) bool {
	if !q.HasMin && !q.HasMax {
		return true
	}
	in := true
	if q.HasMin && width < q.MinWidth {
		in = false
	}
	if q.HasMax && width > q.MaxWidth {
		in = false
	}
	if q.Negated {
		return !in
	}
	return in
}

// matches reports whether any query alternative holds at this width. An empty
// query list (a rule with no @media at all) always matches.
func (m Media) matches(width int) bool {
	if len(m.Queries) == 0 {
		return true
	}
	for _, q := range m.Queries {
		if q.matches(width) {
			return true
		}
	}
	return false
}

// and conjoins two media lists the way nested @media requires: every pair of
// query alternatives intersects its width windows. Either side without queries
// passes the other through unchanged.
func (m Media) and(n Media) Media {
	if len(m.Queries) == 0 {
		return n
	}
	if len(n.Queries) == 0 {
		return m
	}
	var out Media
	for _, a := range m.Queries {
		for _, b := range n.Queries {
			out.Queries = append(out.Queries, intersectQuery(a, b))
		}
	}
	return out
}

// intersectQuery narrows two width windows to their overlap. Negation is left
// on the outer window's side only when it is safe to guess, which in practice
// means: keep whichever side actually constrains the width, and a pair where
// the negations disagree degrades to "always matches" (reported, not fatal).
func intersectQuery(a, b MediaQuery) MediaQuery {
	if a.Negated != b.Negated {
		// A negated window AND a normal one is rare; fall back to the normal
		// side so the rule is not silently dropped.
		if !a.Negated {
			return a
		}
		return b
	}
	q := MediaQuery{Negated: a.Negated}
	if a.HasMin {
		q.MinWidth, q.HasMin = a.MinWidth, true
	}
	if b.HasMin && (!q.HasMin || b.MinWidth > q.MinWidth) {
		q.MinWidth, q.HasMin = b.MinWidth, true
	}
	if a.HasMax {
		q.MaxWidth, q.HasMax = a.MaxWidth, true
	}
	if b.HasMax && (!q.HasMax || b.MaxWidth < q.MaxWidth) {
		q.MaxWidth, q.HasMax = b.MaxWidth, true
	}
	return q
}

// Declaration is a single property/value pair.
type Declaration struct {
	Prop      string
	Raw       string
	Important bool // the declaration carried !important
}

// Rule is one selector-list block, optionally constrained by a media query.
type Rule struct {
	Media     Media
	Selectors []Selector
	Decls     []Declaration
	Order     int
}

// Sheet is a parsed stylesheet: every rule in order, the @keyframes blocks it
// named, plus warnings for the declarations it did not understand. Like a
// browser, the engine skips what it does not know and keeps going; the
// warnings let a template tell its author.
type Sheet struct {
	rules     []*Rule
	order     int
	keyframes map[string]*Keyframes
	Warn      []string
}

// Parse turns CSS source text into a Sheet. Unknown declarations are skipped
// and reported through Warn; a malformed selector or an unbalanced block is an
// error.
func Parse(src string) (*Sheet, error) {
	p := &parser{src: src}
	sh := &Sheet{}
	for {
		p.skipSpace()
		if p.eof() {
			return sh, nil
		}
		switch {
		case p.peek() == '@':
			if e := p.parseAtRule(sh, Media{}); e != nil {
				return sh, e
			}
		case p.peek() == ';':
			p.next()
		default:
			if e := p.parseRule(sh, Media{}); e != nil {
				return sh, e
			}
		}
	}
}

// Rules returns every rule in the sheet, in stylesheet order.
func (sh *Sheet) Rules() []*Rule { return sh.rules }

// Property reads the winning raw value for one property of the element at
// rest, the way the cascade would compute it: the last matching declaration
// in stylesheet order, with any !important declaration overriding a normal one.
// It is how a template answers "what does my button's background say" without
// computing the whole style.
func (sh *Sheet) Property(tag string, classes []string, prop string, baseWidth int) (string, bool) {
	prop = strings.ToLower(strings.TrimSpace(prop))
	var normalRaw string
	var impRaw string
	var hasNormal, hasImp bool
	for _, r := range sh.rules {
		if !r.Media.matches(baseWidth) {
			continue
		}
		match := false
		for _, sel := range r.Selectors {
			if sel.match(tag, classes, StateNone) {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		for _, d := range r.Decls {
			if strings.EqualFold(strings.TrimSpace(d.Prop), prop) {
				raw := strings.TrimSpace(d.Raw)
				if d.Important {
					impRaw, hasImp = raw, true
				} else {
					normalRaw, hasNormal = raw, true
				}
			}
		}
	}
	if hasImp {
		return impRaw, true
	}
	return normalRaw, hasNormal
}

// SetProperty overrides a property for every element with the given classes,
// as a stylesheet appended after all the others — so it wins the cascade
// against anything the sheet said, while media queries still test it. Set
// classes to nil to target the element by tag alone.
func (sh *Sheet) SetProperty(tag string, classes []string, prop, raw string) {
	s := Selector{Tag: tag, Classes: append([]string(nil), classes...)}
	sh.rules = append(sh.rules, &Rule{
		Selectors: []Selector{s},
		Decls:     []Declaration{{Prop: prop, Raw: raw}},
		Order:     sh.order,
	})
	sh.order++
}

// Style computes the winning declarations for an element: matching rules are
// folded by importance, specificity and source order, exactly as a web browser
// would. Custom properties (the --foo kind) are resolved first, then normal
// declarations — and every !important declaration applies last, overwriting
// anything from the same cascade scope. baseWidth is the window width, which
// media queries test against and percentage/viewport lengths scale with. The
// method appends warnings about vendor-prefixed properties, unknown properties
// and unrecognised media conditions to the sheet's warning list.
func (sh *Sheet) Style(tag string, classes []string, state State, baseWidth int) Style {
	return sh.StyleUnits(tag, classes, state, Units{Width: baseWidth, Height: baseWidth})
}

// StyleUnits is [Sheet.Style] with the full measurement context: the window
// size for vw/vh/vmin/vmax and percentages, and the font sizes for em/rem.
// The context's Font is also the cascade hook — a font-size declaration living
// inside the rules updates Font before later em lengths resolve.
func (sh *Sheet) StyleUnits(tag string, classes []string, state State, ctx Units) Style {
	return sh.styleUnits(tag, classes, state, ctx, true)
}

// styleUnits resolves one element's style. inherit folds the body's computed
// value into the element for the inherited properties; the body itself is
// computed with inherit off, so the recursion bottoms out in the theme.
func (sh *Sheet) styleUnits(tag string, classes []string, state State, ctx Units, inherit bool) Style {
	type candidate struct {
		spec  [3]int
		order int
		decls []Declaration
	}
	var got []candidate
	for _, r := range sh.rules {
		if !r.Media.matches(ctx.Width) {
			continue
		}
		for _, sel := range r.Selectors {
			if !sel.match(tag, classes, state) {
				continue
			}
			_, b, c := sel.specificity()
			got = append(got, candidate{[3]int{0, b, c}, r.Order, r.Decls})
			break
		}
	}
	sort.SliceStable(got, func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if got[i].spec[k] != got[j].spec[k] {
				return got[i].spec[k] < got[j].spec[k]
			}
		}
		return false
	})

	// Pass 1: resolve custom properties (--foo: …) in cascade order.
	customs := make(map[string]string, 8)
	for _, cand := range got {
		for _, d := range cand.decls {
			prop := strings.ToLower(strings.TrimSpace(d.Prop))
			if strings.HasPrefix(prop, "--") {
				customs[prop] = d.Raw
			}
		}
	}

	// Partition remaining declarations into normal and !important buckets,
	// each keeping source order.
	var normal, imp []Declaration
	for _, cand := range got {
		for _, d := range cand.decls {
			prop := strings.ToLower(strings.TrimSpace(d.Prop))
			if strings.HasPrefix(prop, "--") {
				continue
			}
			if d.Important {
				imp = append(imp, d)
			} else {
				normal = append(normal, d)
			}
		}
	}

	var st Style
	st.Custom = make(map[string]string, 8)
	st.inherit = make(map[string]bool, 4)
	for p, v := range customs {
		st.Custom[p] = v
	}
	set := map[string]bool{}
	var warns []string
	for _, d := range normal {
		w := applyDecl(&st, set, st.Custom, customs, d, &ctx)
		warns = append(warns, w...)
	}
	for _, d := range imp {
		w := applyDecl(&st, set, st.Custom, customs, d, &ctx)
		warns = append(warns, w...)
	}
	if inherit {
		applyInheritance(&st, set, sh, state, ctx)
	}
	st.Set = set
	resolveCurrentColors(&st)
	finishTransitions(&st)
	finishAnimations(&st)
	sh.Warn = append(sh.Warn, warns...)
	return st
}

// resolveCurrentColors replaces the currentColor sentinel in every colour
// field with the cascade's computed color property. A style that never set
// color leaves the field at zero, which a template reads as its theme
// foreground. color: currentColor itself is a cycle and also falls back to
// zero — the sentinel must never leave the cache.
func resolveCurrentColors(st *Style) {
	if st.Color == CurrentColor {
		st.Color = 0
	}
	if st.Background == CurrentColor {
		st.Background = currentInk(st)
	}
	for i := range st.BoxColor {
		if st.BoxColor[i] == CurrentColor {
			st.BoxColor[i] = currentInk(st)
		}
	}
	for i := range st.BoxShadow {
		if st.BoxShadow[i].Color == CurrentColor {
			st.BoxShadow[i].Color = currentInk(st)
		}
	}
	for i := range st.TextShadow {
		if st.TextShadow[i].Color == CurrentColor {
			st.TextShadow[i].Color = currentInk(st)
		}
	}
	if st.OutlineColor == CurrentColor {
		st.OutlineColor = currentInk(st)
	}
	for i := range st.Filters {
		if st.Filters[i].Drop != nil && st.Filters[i].Drop.Color == CurrentColor {
			st.Filters[i].Drop.Color = currentInk(st)
		}
	}
	for i := range st.BackgroundImages {
		if st.BackgroundImages[i].Grad == nil {
			continue
		}
		for j := range st.BackgroundImages[i].Grad.Stops {
			if st.BackgroundImages[i].Grad.Stops[j].Color == CurrentColor {
				st.BackgroundImages[i].Grad.Stops[j].Color = currentInk(st)
			}
		}
	}
}

func currentInk(st *Style) canvas.Color {
	if st.Set["color"] && st.Color != CurrentColor {
		return st.Color
	}
	return 0
}

// applyInheritance folds the body's computed style into the element in two
// passes: first the properties whose cascade explicitly chose the inherit
// keyword (or unset on an inherited property), then every inherited property
// the element never mentioned. The body copy is lazy and computed once per
// StyleUnits call with inheritance off, so there is no recursion loop.
func applyInheritance(st *Style, set map[string]bool, sh *Sheet, state State, ctx Units) {
	body := sh.styleUnits("body", nil, state, ctx, false)
	for prop := range st.inherit {
		inheritOne(st, set, body, prop)
	}
	for prop := range inheritedProps {
		if !set[prop] && body.Set[prop] {
			inheritOne(st, set, body, prop)
		}
	}
}

// inheritOne copies one computed value from the body style into the element
// and records that the property came by inheritance. The body holds the
// winning cascade value already—including a value the body itself inherited
// from nothing, which is the theme's zero.
func inheritOne(st *Style, set map[string]bool, body Style, prop string) {
	switch prop {
	case "color":
		st.Color = body.Color
	case "font-size":
		st.FontSize = body.FontSize
	case "font-weight":
		st.FontWeight = body.FontWeight
	case "font-style":
		st.FontStyle = body.FontStyle
	case "line-height":
		st.LineHeight = body.LineHeight
	case "letter-spacing":
		st.LetterSpacing = body.LetterSpacing
	case "word-spacing":
		st.WordSpacing = body.WordSpacing
	case "text-align":
		st.TextAlign = body.TextAlign
	case "text-transform":
		st.TextTransform = body.TextTransform
	case "white-space":
		st.WhiteSpace = body.WhiteSpace
	case "overflow-wrap":
		st.OverflowWrap = body.OverflowWrap
	case "visibility":
		st.Visibility = body.Visibility
	case "cursor":
		st.Cursor = body.Cursor
	case "text-shadow":
		st.TextShadow = body.TextShadow
	}
	set[prop] = true
	delete(st.inherit, prop)
}

func fmtErrf(format string, args ...any) error {
	return fmt.Errorf("css: "+format, args...)
}

func splitWords(raw string) []string {
	return strings.Fields(raw)
}
