package css

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
	if st.ColumnRuleColor == CurrentColor {
		st.ColumnRuleColor = currentInk(st)
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
