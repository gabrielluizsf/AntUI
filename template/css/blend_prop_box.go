package css

// blendPropBox interpolates the geometric properties, which lerp their
// lengths against the ruler the Units carry.
func blendPropBox(dst *Style, from Style, prop string, t float64, ctx Units) {
	switch prop {
	case "width":
		dst.Width = blendLength(from.Width, dst.Width, t, ctx)
	case "height":
		dst.Height = blendLength(from.Height, dst.Height, t, ctx)
	case "min-width":
		dst.MinWidth = blendLength(from.MinWidth, dst.MinWidth, t, ctx)
	case "max-width":
		dst.MaxWidth = blendLength(from.MaxWidth, dst.MaxWidth, t, ctx)
	case "min-height":
		dst.MinHeight = blendLength(from.MinHeight, dst.MinHeight, t, ctx)
	case "max-height":
		dst.MaxHeight = blendLength(from.MaxHeight, dst.MaxHeight, t, ctx)
	case "margin-top":
		dst.Margin[0] = blendLength(from.Margin[0], dst.Margin[0], t, ctx)
	case "margin-right":
		dst.Margin[1] = blendLength(from.Margin[1], dst.Margin[1], t, ctx)
	case "margin-bottom":
		dst.Margin[2] = blendLength(from.Margin[2], dst.Margin[2], t, ctx)
	case "margin-left":
		dst.Margin[3] = blendLength(from.Margin[3], dst.Margin[3], t, ctx)
	case "padding-top":
		dst.Padding[0] = blendLength(from.Padding[0], dst.Padding[0], t, ctx)
	case "padding-right":
		dst.Padding[1] = blendLength(from.Padding[1], dst.Padding[1], t, ctx)
	case "padding-bottom":
		dst.Padding[2] = blendLength(from.Padding[2], dst.Padding[2], t, ctx)
	case "padding-left":
		dst.Padding[3] = blendLength(from.Padding[3], dst.Padding[3], t, ctx)
	case "top":
		dst.Top = blendLength(from.Top, dst.Top, t, ctx)
	case "right":
		dst.Right = blendLength(from.Right, dst.Right, t, ctx)
	case "bottom":
		dst.Bottom = blendLength(from.Bottom, dst.Bottom, t, ctx)
	case "left":
		dst.Left = blendLength(from.Left, dst.Left, t, ctx)
	case "font-size":
		dst.FontSize = int(float64(from.FontSize) + float64(dst.FontSize-from.FontSize)*t + 0.5)
	case "letter-spacing":
		dst.LetterSpacing = blendLength(from.LetterSpacing, dst.LetterSpacing, t, ctx)
	case "word-spacing":
		dst.WordSpacing = blendLength(from.WordSpacing, dst.WordSpacing, t, ctx)
	case "line-height":
		dst.LineHeight = from.LineHeight + (dst.LineHeight-from.LineHeight)*t
	}
}
