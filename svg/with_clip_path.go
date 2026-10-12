package svg

import "strings"

// withClipPath keeps the id a `clip-path` names: the shapes that cut this
// element are elsewhere in the drawing, the same way a gradient is, so only a
// caller holding the whole drawing can say what it is. See
// [Image.resolveClip].
func (s *Style) withClipPath(raw string, warn func(string, ...any)) {
	switch id, ok := paintRef(raw); {
	case strings.EqualFold(strings.TrimSpace(raw), "none"):
		s.clipRef = ""
	case ok:
		s.clipRef = id
	default:
		s.clipRef = ""
		warn("the clip-path %q is not a reference to a <clipPath> this package can follow, so the element is drawn without it", raw)
	}
}

// withMask keeps the id a `mask` names, the same way [Style.withClipPath]
// keeps a clip's: only the whole drawing can say what it is. It is not
// inherited — see the `maskRef` field of [Style] and [Image.maskNamed].
func (s *Style) withMask(raw string, warn func(string, ...any)) {
	switch id, ok := paintRef(raw); {
	case strings.EqualFold(strings.TrimSpace(raw), "none"):
		s.maskRef = ""
	case ok:
		s.maskRef = id
	default:
		s.maskRef = ""
		warn("the mask %q is not a reference to a <mask> this package can follow, so the element is drawn without it", raw)
	}
}
