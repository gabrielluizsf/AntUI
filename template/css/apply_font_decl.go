package css

func applyFontDecl(e *declEnv) []string {
	st, raw, ctx, note := e.st, e.raw, e.ctx, e.note
	switch e.prop {
	case "font-size":
		if l, err := parseLengthAt(raw, *ctx); err == nil && !l.IsPct() {
			st.FontSize = l.Resolve(*ctx)
			if st.FontSize <= 0 {
				st.FontSize = DefaultFontSize
			}
			ctx.Font = st.FontSize
			note()
		}
	case "font-weight":
		if w, ok := parseFontWeight(raw); ok {
			st.FontWeight = w
			note()
		}
	case "font-style":
		if v, ok := parseFontStyle(raw); ok {
			st.FontStyle = v
			note()
		}
	case "font-family":
		if v := normalizeFamily(raw); v != "" {
			st.FontFamily = v
			note()
		}
	case "line-height":
		if v, ok := parseLineHeight(raw, *ctx); ok {
			st.LineHeight = v
			note()
		}
	case "letter-spacing":
		if raw == "normal" {
			st.LetterSpacing = Zero()
			note()
		} else if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.LetterSpacing = l
			note()
		}
	case "word-spacing":
		if raw == "normal" {
			st.WordSpacing = Zero()
			note()
		} else if l, err := parseLengthAt(raw, *ctx); err == nil {
			st.WordSpacing = l
			note()
		}
	}
	return nil
}
