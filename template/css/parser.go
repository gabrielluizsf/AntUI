package css

import (
	"fmt"
	"strconv"
	"strings"
)

// parser walks CSS source text. It is deliberately forgiving: a selector,
// at-rule or property the engine does not understand is dropped with a
// warning, and only structural damage — an unterminated block, a rule with no
// opening brace — is an error. The tokenizer also never spins: every scan step
// either advances the index or stops, so a stylesheet of any kind can spend
// forever being read rather than hang the app.
type parser struct {
	src string
	i   int

	// path is the file this text came from, cleaned, and empty when it came
	// from no file at all — which is also what leaves a relative @import
	// with nowhere to resolve to.
	path string

	// seen is every file this parse has read or is reading. An @import of
	// one of them would read it a second time, and two files importing each
	// other would never end, so the whole tree of files shares one list.
	seen map[string]bool

	// depth is how many @imports deep this text was read — the sheet's own
	// file is 0. It bounds a chain in which every file is a different one
	// and the seen list alone would let it run on.
	depth int
}

// maxImportDepth is how many @imports may nest: one sheet reading another,
// that reading another, and so on. A chain longer than this is reported and
// left unread rather than followed to wherever it goes.
const maxImportDepth = 16

func (p *parser) eof() bool { return p.i >= len(p.src) }

func (p *parser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.src[p.i]
}

func (p *parser) next() byte {
	c := p.src[p.i]
	p.i++
	return c
}

// parseMedia reads an @media prelude into a Media. The full grammar is
// walked — not/only, and/or, comma-separated query lists — while the
// features the viewport answers constrain the rule (its two edges, the
// window's shape, the display's density and the system's color scheme);
// everything else warns and is treated as satisfied, so a canvas never drops
// a rule it merely cannot measure.
func parseMedia(prelude string) (Media, []string) {
	var m Media
	var warns []string
	for _, branch := range splitTopLevel(prelude) {
		qs, w := parseMediaQuery(branch)
		warns = append(warns, w...)
		m.Queries = append(m.Queries, qs...)
	}
	return m, warns
}

// parseMediaQuery parses one comma-separated media query into its alternative
// width windows, or-separated groups each becoming a query of their own.
func parseMediaQuery(branch string) ([]MediaQuery, []string) {
	var out []MediaQuery
	var warns []string
	negated := false
	cur := MediaQuery{}
	flush := func() {
		if !cur.constrained() {
			// A negation of nothing we can evaluate is left unconstrained so
			// the rule still applies.
			cur.Negated = false
		}
		out = append(out, cur)
	}
	for _, tok := range splitQueryTokens(branch) {
		t := strings.ToLower(strings.TrimSpace(tok))
		if t == "" {
			continue
		}
		switch t {
		case "not":
			negated = true
			cur.Negated = true
		case "only":
			// A media-type qualifier; the type itself is beyond the canvas.
		case "and":
			// Conjunction: keep filling the current query.
		case "or":
			flush()
			cur = MediaQuery{Negated: negated}
		default:
			if t[0] == '(' {
				body := strings.TrimSuffix(strings.TrimPrefix(t, "("), ")")
				kv := strings.SplitN(body, ":", 2)
				if len(kv) != 2 {
					warns = append(warns, fmtErrf("ignoring media condition %q", t).Error())
					continue
				}
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				n, ok := mediaPx(val)
				// measure writes one edge of the window, or says the number
				// was not a number and leaves the edge unset.
				measure := func(dst *int, has *bool, edge string) {
					if ok {
						*dst, *has = n, true
						return
					}
					warns = append(warns, fmtErrf("ignoring media %s %q", edge, val).Error())
				}
				switch key {
				case "min-width":
					measure(&cur.MinWidth, &cur.HasMin, "width")
				case "max-width":
					measure(&cur.MaxWidth, &cur.HasMax, "width")
				case "min-height":
					measure(&cur.MinHeight, &cur.HasMinHeight, "height")
				case "max-height":
					measure(&cur.MaxHeight, &cur.HasMaxHeight, "height")
				case "orientation":
					switch val {
					case "portrait":
						cur.HasOrientation, cur.Portrait = true, true
					case "landscape":
						cur.HasOrientation, cur.Portrait = true, false
					default:
						warns = append(warns, fmtErrf("ignoring media orientation %q", val).Error())
					}
				case "resolution", "min-resolution", "max-resolution":
					dpi, ok := mediaDpi(val)
					if !ok {
						warns = append(warns, fmtErrf("ignoring media resolution %q", val).Error())
						continue
					}
					switch key {
					case "resolution":
						cur.MinDpi, cur.MaxDpi, cur.HasMinDpi, cur.HasMaxDpi = dpi, dpi, true, true
					case "min-resolution":
						cur.MinDpi, cur.HasMinDpi = dpi, true
					case "max-resolution":
						cur.MaxDpi, cur.HasMaxDpi = dpi, true
					}
				case "prefers-color-scheme":
					switch val {
					case "light":
						cur.HasScheme, cur.Dark = true, false
					case "dark":
						cur.HasScheme, cur.Dark = true, true
					default:
						warns = append(warns, fmtErrf("ignoring media color-scheme %q", val).Error())
					}
				default:
					warns = append(warns, fmtErrf("ignoring media condition %q", t).Error())
				}
			} else {
				switch t {
				case "screen", "all":
					// The window kind the canvas always is.
				default:
					warns = append(warns, fmtErrf("ignoring media type %q", t).Error())
				}
			}
		}
	}
	flush()
	return out, warns
}

// splitQueryTokens breaks a media query into its words and (…)-groups,
// keeping parenthesised groups — including nested parentheses — whole.
func splitQueryTokens(s string) []string {
	var out []string
	i := 0
	for i < len(s) {
		for i < len(s) && isSpace(s[i]) {
			i++
		}
		if i >= len(s) {
			break
		}
		start := i
		if s[i] == '(' {
			i = max(scanParen(s, i), i+1)
		} else {
			for i < len(s) && !isSpace(s[i]) && s[i] != '(' {
				i++
			}
		}
		if i == start {
			i++
			continue
		}
		out = append(out, s[start:i])
	}
	return out
}

// mediaPx reads a media width like "720px" or "720" into pixels. Fractions
// round; junk returns ok=false and the callers warn.
func mediaPx(s string) (int, bool) {
	t := strings.TrimSpace(strings.ToLower(strings.TrimSuffix(strings.TrimSpace(s), "px")))
	var f float64
	if _, err := fmt.Sscanf(t, "%g", &f); err != nil {
		return 0, false
	}
	return int(f + 0.5), true
}

// mediaDpi reads a resolution — 96dpi, 40dpcm, 2dppx — and answers it in dots
// per inch, the unit every query is compared in. A bare number, a value with
// no digits in it, or a unit the spec does not give a resolution returns
// ok=false and the caller warns.
func mediaDpi(s string) (float64, bool) {
	t := strings.TrimSpace(strings.ToLower(s))
	i := 0
	for i < len(t) && (t[i] == '.' || (t[i] >= '0' && t[i] <= '9')) {
		i++
	}
	v, err := strconv.ParseFloat(t[:i], 64)
	if err != nil {
		return 0, false
	}
	switch strings.TrimSpace(t[i:]) {
	case "dpi":
	case "dpcm":
		v *= 2.54
	case "dppx":
		v *= 96
	default:
		return 0, false
	}
	return roundDpi(v), true
}

// stripComments removes every /* … */ comment from a string. A comment in a
// file any old position — selector, value, media prelude — carries nothing, so
// it is just deleted.
func stripComments(s string) string {
	for {
		a := strings.Index(s, "/*")
		if a < 0 {
			return s
		}
		b := strings.Index(s[a+2:], "*/")
		if b < 0 {
			return s[:a]
		}
		end := a + 2 + b + 2
		s = s[:a] + s[end:]
	}
}

// parseDeclaration reads "property: value". A missing ':' is not a
// declaration. A trailing !important is recorded, not parsed as a value.
func parseDeclaration(text string) (Declaration, bool) {
	text = stripComments(text)
	idx := strings.IndexByte(text, ':')
	if idx < 0 {
		return Declaration{}, false
	}
	prop := strings.TrimSpace(text[:idx])
	val := strings.TrimSpace(text[idx+1:])
	if prop == "" {
		return Declaration{}, false
	}
	val, important := splitImportant(val)
	return Declaration{Prop: prop, Raw: val, Important: important}, true
}

// splitImportant peels a trailing "!important" (in any spacing) off a value.
func splitImportant(raw string) (val string, important bool) {
	v := strings.TrimSpace(raw)
	if i := strings.LastIndexByte(v, '!'); i >= 0 {
		rest := strings.TrimSpace(strings.ToLower(v[i+1:]))
		if rest == "important" {
			return strings.TrimSpace(v[:i]), true
		}
	}
	return v, false
}

// scanFloat reads a floating point number with fmt-equivalent strictness.
func scanFloat(s string, dst *float64) (int, error) {
	t := strings.TrimSuffix(s, "%")
	v, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
	if err != nil {
		return 0, err
	}
	*dst = v
	if strings.HasSuffix(strings.TrimSpace(s), "%") {
		*dst = v / 100
	}
	return len(t), nil
}

// isSpace reports whether c is CSS whitespace.
func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '\f':
		return true
	}
	return false
}

func isIdentByte(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
