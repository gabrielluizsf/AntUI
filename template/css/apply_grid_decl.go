package css

func applyGridDecl(e *declEnv) []string {
	st, raw, ctx, note := e.st, e.raw, e.ctx, e.note
	switch e.prop {
	case "grid-template-columns":
		if v, ok := parseGridTemplate(raw, *ctx); ok {
			st.GridTemplateColumns = v
			note()
		}
	case "grid-template-rows":
		if v, ok := parseGridTemplate(raw, *ctx); ok {
			st.GridTemplateRows = v
			note()
		}
	case "grid-template-areas":
		if v, ok := parseGridTemplateAreas(raw); ok {
			st.GridTemplateAreas = v
			note()
		}
	case "grid-auto-flow":
		if flow, dense, ok := parseGridAutoFlow(raw); ok {
			st.GridAutoFlow = flow
			st.GridAutoFlowDense = dense
			note()
		}
	case "justify-items":
		if v, ok := parseJustifyItems(raw); ok {
			st.JustifyItems = v
			note()
		}
	case "justify-self":
		if v, ok := parseJustifySelf(raw); ok {
			st.JustifySelf = v
			note()
		}
	}
	return nil
}
