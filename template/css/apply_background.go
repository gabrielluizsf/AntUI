package css

import (
	"strings"
)

// applyBackground folds a background shorthand into the style: a comma
// separates layers, position and size split at "/", and a colour may only be
// the last layer's. A value that does not fit the grammar drops the whole
// declaration, the way an unrecognised value does.
func applyBackground(st *Style, raw string, ctx Units) bool {
	layers := splitFields(raw, ',')
	if len(layers) == 0 || layers[len(layers)-1] == "" {
		return false
	}
	if len(layers) == 1 && strings.TrimSpace(layers[0]) == "none" {
		st.BackgroundImages = nil
		return true
	}
	var fold backgroundFold
	for idx, layerStr := range layers {
		part := strings.TrimSpace(layerStr)
		head, sizeStr, hasSlash := splitSlash(part)
		var size BackSize
		if hasSlash {
			var ok bool
			if size, ok = parseBackSizeOne(strings.TrimSpace(sizeStr), ctx); !ok {
				return false
			}
		}
		if !fold.add(scanBackgroundLayer(head, idx == len(layers)-1, ctx), size, ctx) {
			return false
		}
	}
	fold.store(st)
	return true
}
