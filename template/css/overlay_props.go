package css

// overlayProps folds every property in props from src into dst, marking each
// in dst's Set map.
func overlayProps(dst *Style, src Style, props map[string]bool) {
	for prop := range props {
		setProp(dst, src, prop)
	}
}
