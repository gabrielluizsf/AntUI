package css

// applyMediaResolution reads a resolution feature — 96dpi, 40dpcm, 2dppx —
// and folds it into the query's one or two dpi edges.
func applyMediaResolution(cur *MediaQuery, f mediaFeat, warns []string) []string {
	dpi, ok := mediaDpi(f.val)
	if !ok {
		return append(warns, fmtErrf("ignoring media resolution %q", f.val).Error())
	}
	switch f.key {
	case "resolution":
		cur.MinDpi, cur.MaxDpi, cur.HasMinDpi, cur.HasMaxDpi = dpi, dpi, true, true
	case "min-resolution":
		cur.MinDpi, cur.HasMinDpi = dpi, true
	case "max-resolution":
		cur.MaxDpi, cur.HasMaxDpi = dpi, true
	}
	return warns
}
