package css

// mediaFeat is one parenthesised media feature: its name, its value, and the
// raw token it came from, kept for the warning an unreadable one produces.
type mediaFeat struct {
	key string
	val string
	tok string
}

// applyMediaFeature folds one media feature into the query, updating the edge
// it constrains and returning the warnings for the ones it cannot read.
func applyMediaFeature(cur *MediaQuery, f mediaFeat) []string {
	n, ok := mediaPx(f.val)
	var warns []string
	// measure writes one edge of the window, or says the number was not a
	// number and leaves the edge unset.
	measure := func(dst *int, has *bool, edge string) {
		if ok {
			*dst, *has = n, true
			return
		}
		warns = append(warns, fmtErrf("ignoring media %s %q", edge, f.val).Error())
	}
	switch f.key {
	case "min-width":
		measure(&cur.MinWidth, &cur.HasMin, "width")
	case "max-width":
		measure(&cur.MaxWidth, &cur.HasMax, "width")
	case "min-height":
		measure(&cur.MinHeight, &cur.HasMinHeight, "height")
	case "max-height":
		measure(&cur.MaxHeight, &cur.HasMaxHeight, "height")
	case "orientation":
		switch f.val {
		case "portrait":
			cur.HasOrientation, cur.Portrait = true, true
		case "landscape":
			cur.HasOrientation, cur.Portrait = true, false
		default:
			warns = append(warns, fmtErrf("ignoring media orientation %q", f.val).Error())
		}
	case "resolution", "min-resolution", "max-resolution":
		return applyMediaResolution(cur, f, warns)
	case "prefers-color-scheme":
		switch f.val {
		case "light":
			cur.HasScheme, cur.Dark = true, false
		case "dark":
			cur.HasScheme, cur.Dark = true, true
		default:
			warns = append(warns, fmtErrf("ignoring media color-scheme %q", f.val).Error())
		}
	default:
		warns = append(warns, fmtErrf("ignoring media condition %q", f.tok).Error())
	}
	return warns
}
