package svg

import "github.com/gabrielluizsf/antui/canvas"

// nonScalingDashed cuts a non-scaling stroke after its path has reached the
// device. Dash lengths and offsets are device pixels, just like its width;
// supersampling keeps the cut at the same final-pixel size as the outline.
func nonScalingDashed(n *Node, m canvas.Matrix) *canvas.Path {
	source := n.Style.Dash
	if source == nil {
		return nil
	}
	lengths := make([]float64, 0, len(source.On)+len(source.Off))
	for _, length := range source.On {
		lengths = append(lengths, length*supersample)
	}
	for _, length := range source.Off {
		lengths = append(lengths, length*supersample)
	}
	dash, ok := canvas.NewDash(lengths, source.Offset*supersample)
	if !ok {
		return nil
	}

	path := clonePath(n.Path)
	path.Transform(m)
	return path.Dashed(dash)
}
