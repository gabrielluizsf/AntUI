package svg

// sheetDecls is what the drawing's own `<style>` sheets say about one element,
// normal and important kept apart because an important declaration beats every
// normal one, whether it came from a sheet or from the element's own `style`.
type sheetDecls struct {
	normal    map[string]string
	important map[string]string
}

// sheetMap turns the winning declarations into the element's sheet style,
// splitting them by importance.
func sheetMap(best map[string]cssWinner) *sheetDecls {
	var out *sheetDecls
	for name, w := range best {
		if out == nil {
			out = &sheetDecls{normal: map[string]string{}, important: map[string]string{}}
		}
		if w.important {
			out.important[name] = w.value
		} else {
			out.normal[name] = w.value
		}
	}
	return out
}
