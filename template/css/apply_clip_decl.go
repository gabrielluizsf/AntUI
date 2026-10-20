package css

func applyClipDecl(e *declEnv) []string {
	if e.prop != "clip-path" {
		return nil
	}
	p, ok := parseClipPath(e.raw)
	if !ok {
		return []string{fmtErrf("ignoring clip-path %q: only polygon() is supported", e.raw).Error()}
	}
	e.st.ClipPath = p
	e.note()
	return nil
}
