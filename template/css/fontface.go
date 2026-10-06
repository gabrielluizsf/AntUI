package css

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// FontFace is one @font-face block: the family, weight and slant it was
// declared for, and the file it read. A face is read while the sheet parses,
// so a frame only ever looks one up.
type FontFace struct {
	Family string // lowercased, the way font-family names match
	Weight uint16 // the weight it declared, on the hundreds CSS matches on
	Style  uint8  // one of the FontStyle* constants
	face   *canvas.Face
}

// Face is the font file itself, at the size the program draws its text at,
// with the names after it in the font-family list behind it for the runes it
// has no shape for.
func (ff *FontFace) Face() *canvas.Face { return ff.face }

// Slanted reports whether the file already leans, so a style that asked for
// italic has no synthetic slant to add on top of it.
func (ff *FontFace) Slanted() bool { return ff.Style != FontStyleNormal }

// Bold reports whether the file already carries the heavy weight, so a style
// that asked for bold has no fake bold to draw over it.
func (ff *FontFace) Bold() bool { return ff.Weight >= 600 }

// fontKey is one question the sheet answers: a font-family list at a weight
// and a slant.
type fontKey struct {
	family  string
	weight  uint16
	slanted bool
}

// Font is the face a font-family list asks for at a weight and a slant, or
// nil when no @font-face in the sheet holds any of its names — the caller
// keeps drawing in the face it already had.
//
// The list is read name by name: the first name the sheet holds a face for
// draws, and every name after it is chained behind that face so the runes it
// lacks reach the next one that has them, ending at the face the program
// draws with by default. Inside one name CSS chooses: a face already leaning
// is looked for first when the style leans, and a straight one first when it
// does not; within that, the wanted weight and then the ones the spec tries
// from it.
//
// The answer is kept, so a frame asks once however many lines it draws — and
// so the chain it hands back ends at the default face the program had the
// first time the question was asked.
func (sh *Sheet) Font(family string, weight uint16, slanted bool) *FontFace {
	if sh == nil || family == "" {
		return nil
	}
	key := fontKey{family: family, weight: normWeight(weight), slanted: slanted}
	sh.fontMu.RLock()
	ff, seen := sh.fontCache[key]
	sh.fontMu.RUnlock()
	if seen {
		return ff
	}
	ff = sh.fontChain(family, weight, slanted)
	sh.fontMu.Lock()
	if sh.fontCache == nil {
		sh.fontCache = map[fontKey]*FontFace{}
	}
	sh.fontCache[key] = ff
	sh.fontMu.Unlock()
	return ff
}

// fontChain is one answer to a Font query, with the faces of the names the
// list holds chained behind each other, or nil when it holds none of them.
func (sh *Sheet) fontChain(family string, weight uint16, slanted bool) *FontFace {
	var chain []*FontFace
	for _, name := range splitTopLevel(family) {
		if ff := sh.fontNamed(unquote(name), weight, slanted); ff != nil {
			chain = append(chain, ff)
		}
	}
	if len(chain) == 0 {
		return nil
	}
	// Every face in the chain is drawn at the size the program draws the
	// rest of its text at, so naming a family does not change how big the
	// text comes out — the stylesheet's own font-size is what sizes it.
	base := canvas.DefaultFace().Size()
	if base <= 0 {
		base = float64(DefaultFontSize)
	}
	// The chain is built from the end back: each face takes the one after
	// it as its fallback, so the list runs out at the face the program
	// draws with by default.
	face := canvas.DefaultFace()
	for i := len(chain) - 1; i >= 0; i-- {
		face = chain[i].face.AtSize(base).WithFallback(face)
	}
	head := *chain[0]
	head.face = face
	return &head
}

// fontNamed is the face the sheet holds for one family name at a weight and a
// slant, or nil when it holds none.
func (sh *Sheet) fontNamed(name string, weight uint16, slanted bool) *FontFace {
	want := normWeight(weight)
	for _, lean := range [2]bool{slanted, !slanted} {
		var pick *FontFace
		best := -1
		for _, ff := range sh.fonts {
			if !strings.EqualFold(ff.Family, name) || ff.Slanted() != lean {
				continue
			}
			if r := weightRank(want, ff.Weight); pick == nil || r < best {
				pick, best = ff, r
			}
		}
		if pick != nil {
			return pick
		}
	}
	return nil
}

// weightRank is how soon CSS Fonts asks for w when want is the weight in
// hand: 0 is the exact weight, and the rest follow the order the search
// turns — down first below 500, up first at 500 or above, with 500 ahead of
// every weight below 400 when 400 is asked for.
func weightRank(want, w uint16) int {
	w = normWeight(w)
	if w == want {
		return 0
	}
	switch {
	case want == 400:
		if w == 500 {
			return 1
		}
		if w < want {
			return 1 + int((want-w)/100)
		}
		// Past the weights below 400, the climb starts at 600: 500 has
		// its own place already.
		return 4 + int((w-500)/100)
	case want < 400:
		if w < want {
			return int((want - w) / 100)
		}
		return int((want-100)/100) + int((w-want)/100)
	case want == 500:
		if w < want {
			return int((want - w) / 100)
		}
		return 5 + int((w-600)/100)
	default:
		if w > want {
			return int((w - want) / 100)
		}
		return int((maxWeight-want)/100) + int((want-w)/100)
	}
}

// maxWeight is the heaviest weight the spec's scale reaches, and normWeight
// rounds a declared weight onto the hundreds the match runs on.
const maxWeight = 1000

func normWeight(w uint16) uint16 {
	if w == 0 {
		return FontWeightNormal
	}
	return min(max((w+50)/100*100, 100), maxWeight)
}

// parseFontFace reads a @font-face block: the family the file answers to,
// which file it is, and the weight and slant it was cut for. Everything the
// block says that this engine cannot use is reported and passed over, and a
// file that will not open leaves the family with no face at all rather than
// failing the stylesheet that asked for it.
func (p *parser) parseFontFace(sh *Sheet) error {
	prelude, _, err := p.readHeader(false)
	if err != nil {
		return err
	}
	if p.peek() == '{' {
		p.next()
	}
	if s := strings.TrimSpace(prelude); s != "" {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %q before a @font-face block", s).Error())
	}

	var family, src, weight, style string
	err = p.declBlock(func(text string) error {
		d, ok := parseDeclaration(text)
		if !ok {
			return nil
		}
		switch prop := strings.ToLower(strings.TrimSpace(d.Prop)); prop {
		case "font-family":
			family = d.Raw
		case "src":
			src = d.Raw
		case "font-weight":
			weight = d.Raw
		case "font-style":
			style = d.Raw
		case "font-display":
			// How long a stylesheet waits for a face it did not read is
			// answered here: this canvas reads the file now or never, so
			// there is nothing for the hint to time.
		default:
			sh.Warn = append(sh.Warn, fmtErrf("ignoring %q in @font-face", prop).Error())
		}
		return nil
	})
	if err != nil {
		return err
	}

	family = normalizeFamily(family)
	if family == "" {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring a @font-face with no font-family").Error())
		return nil
	}
	if strings.TrimSpace(src) == "" {
		sh.Warn = append(sh.Warn, fmtErrf("ignoring the @font-face for %q: no src", family).Error())
		return nil
	}

	ff := &FontFace{Family: family, Weight: FontWeightNormal}
	if s := strings.TrimSpace(weight); s != "" {
		w, ok := parseFontWeight(s)
		if !ok {
			sh.Warn = append(sh.Warn, fmtErrf("ignoring font-weight %q in @font-face for %q", weight, family).Error())
		} else {
			ff.Weight = w
		}
	}
	if s := strings.TrimSpace(style); s != "" {
		v, ok := parseFontStyle(s)
		if !ok {
			sh.Warn = append(sh.Warn, fmtErrf("ignoring font-style %q in @font-face for %q", style, family).Error())
		} else {
			ff.Style = v
		}
	}
	if ff.face = p.fontSource(sh, family, src); ff.face == nil {
		return nil
	}
	sh.fonts = append(sh.fonts, ff)
	sh.fontMu.Lock()
	sh.fontCache = nil
	sh.fontMu.Unlock()
	return nil
}

// fontSource walks an @font-face src list the way a browser does: the first
// source this canvas can open is the face. local() names a font this machine
// already has, and a woff file is a compressed container this engine does not
// unpack — both are reported and passed over, and so is a source that opens
// as no font at all.
func (p *parser) fontSource(sh *Sheet, family, src string) *canvas.Face {
	for _, part := range splitTopLevel(src) {
		kind, arg := sourceKind(part)
		what := fmt.Sprintf("%q in @font-face for %q", arg, family)
		switch kind {
		case "local":
			sh.Warn = append(sh.Warn, fmtErrf("ignoring local() %q in @font-face for %q: a font this machine already has is not read", arg, family).Error())
			continue
		case "url":
		default:
			sh.Warn = append(sh.Warn, fmtErrf("ignoring %q in @font-face for %q: not a url() or local()", part, family).Error())
			continue
		}
		if format := sourceFormat(part); format == "woff" || format == "woff2" {
			sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: woff is not read", what).Error())
			continue
		}
		path, warn := p.localFile(arg, what)
		if warn != "" {
			sh.Warn = append(sh.Warn, warn)
			continue
		}
		data, err := os.ReadFile(path)
		if err == nil {
			var face *canvas.Face
			if face, err = canvas.ParseFace(data, float64(DefaultFontSize)); err == nil {
				return face
			}
		}
		sh.Warn = append(sh.Warn, fmtErrf("ignoring %s: %v", what, err).Error())
	}
	return nil
}

// sourceKind reads one @font-face src source: which function it is and what
// it names, or "" when the source is neither url() nor local().
func sourceKind(part string) (kind, arg string) {
	s := strings.TrimSpace(part)
	for _, k := range []string{"url", "local"} {
		if len(s) > len(k)+1 && strings.EqualFold(s[:len(k)], k) && s[len(k)] == '(' {
			end := parenEnd(s, len(k))
			if end < 0 {
				return "", ""
			}
			return strings.ToLower(k), unquote(strings.TrimSpace(s[len(k)+1 : end]))
		}
	}
	return "", ""
}

// sourceFormat reads the format("…") clause of a src source, lowercased, or
// "" when the source carries none. It is a hint only: a format this engine
// reads is believed, and one it does not is the only reason to pass a source
// over without opening it.
func sourceFormat(part string) string {
	lower := strings.ToLower(part)
	i := strings.Index(lower, "format(")
	if i < 0 {
		return ""
	}
	s := strings.TrimSpace(lower[i+len("format("):])
	end := strings.IndexByte(s, ')')
	if end < 0 {
		return ""
	}
	return unquote(strings.TrimSpace(s[:end]))
}

// normalizeFamily turns a font-family declaration into the name a
// @font-face answers to: comma-separated, each name unquoted, lowercased —
// a family's letters are its letters whatever case the stylesheet wrote.
func normalizeFamily(raw string) string {
	var out []string
	for _, name := range splitTopLevel(raw) {
		if n := strings.ToLower(unquote(name)); n != "" {
			out = append(out, n)
		}
	}
	return strings.Join(out, ", ")
}

// localFile resolves a URL from the file being read to a path on this
// machine, or returns the warning that says why it is not one. what names
// where the URL came from for that warning.
func (p *parser) localFile(ref, what string) (string, string) {
	if isRemoteRef(ref) {
		return "", fmtErrf("ignoring %s: only local files are read", what).Error()
	}
	// A query and a fragment are for a server to answer; the file beside
	// this one is the file beside this one.
	if i := strings.IndexAny(ref, "?#"); i >= 0 {
		ref = ref[:i]
	}
	if strings.TrimSpace(ref) == "" {
		return "", fmtErrf("ignoring %s: no file name", what).Error()
	}
	if filepath.IsAbs(ref) {
		return filepath.Clean(ref), ""
	}
	if p.path == "" {
		return "", fmtErrf("ignoring %s: no file to resolve it against", what).Error()
	}
	return filepath.Join(filepath.Dir(p.path), ref), ""
}
