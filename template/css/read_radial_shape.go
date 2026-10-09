package css

// readRadialShape applies a leading circle/ellipse keyword and the radial size
// keyword, returning the index of the first token that was neither. The
// defaults were seeded by the caller.
func readRadialShape(tokens []string, g *Gradient) int {
	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case "circle":
			g.Shape = GradientCircle
		case "ellipse":
			g.Shape = GradientEllipse
		case "closest-side":
			g.Size = GradientClosestSide
		case "farthest-side":
			g.Size = GradientFarthestSide
		case "closest-corner":
			g.Size = GradientClosestCorner
		case "farthest-corner":
			g.Size = GradientFarthestCorner
		default:
			return i
		}
		i++
	}
	return i
}
