package css

func applyEffectDecl(e *declEnv) []string {
	st, raw, ctx, note := e.st, e.raw, e.ctx, e.note
	switch e.prop {
	case "box-shadow":
		if sh, ok := parseShadows(raw, *ctx, true, 4); ok {
			st.BoxShadow = sh
			note()
		}
	case "text-shadow":
		if sh, ok := parseShadows(raw, *ctx, false, 3); ok {
			st.TextShadow = sh
			note()
		}
	case "filter":
		if f, ok := parseFilters(raw, *ctx); ok {
			st.Filters = f
			note()
		}
	case "backdrop-filter":
		if f, ok := parseFilters(raw, *ctx); ok {
			st.BackdropFilters = f
			note()
		}
	case "transform":
		if f, ok := parseTransform(raw, *ctx); ok {
			st.Transform = f
			note()
		}
	case "transform-origin":
		if o, ok := parseTransformOrigin(raw, *ctx); ok {
			st.TransformOrigin = o
			note()
		}
	}
	return nil
}
