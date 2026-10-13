package svg

import "github.com/gabrielluizsf/antui/canvas"

// cloneDash copies a dash without sharing its slices. NewDash rebuilds the
// canvas' private cycle, which is needed when a pattern was made by hand or
// was scaled by pathLength rather than parsed directly.
func cloneDash(src *canvas.Dash) *canvas.Dash {
	if src == nil {
		return nil
	}
	lengths := make([]float64, 0, len(src.On)+len(src.Off))
	lengths = append(lengths, src.On...)
	lengths = append(lengths, src.Off...)
	dash, ok := canvas.NewDash(lengths, src.Offset)
	if !ok {
		return &canvas.Dash{
			On: append([]float64{}, src.On...), Off: append([]float64{}, src.Off...), Offset: src.Offset,
		}
	}
	dash.Offset = src.Offset
	return &dash
}

// sameDash compares the public parts of a dash. The cycle cache is an
// implementation detail of canvas and does not change what the pattern draws.
func sameDash(a, b *canvas.Dash) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Offset != b.Offset || len(a.On) != len(b.On) || len(a.Off) != len(b.Off) {
		return false
	}
	for i := range a.On {
		if a.On[i] != b.On[i] {
			return false
		}
	}
	for i := range a.Off {
		if a.Off[i] != b.Off[i] {
			return false
		}
	}
	return true
}
