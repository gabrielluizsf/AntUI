package css

// setPropBox copies one geometric property's value from src to dst: the
// element's box, its four sides and its font metrics.
func setPropBox(dst *Style, src Style, prop string) {
	switch prop {
	case "width":
		dst.Width = src.Width
	case "height":
		dst.Height = src.Height
	case "min-width":
		dst.MinWidth = src.MinWidth
	case "max-width":
		dst.MaxWidth = src.MaxWidth
	case "min-height":
		dst.MinHeight = src.MinHeight
	case "max-height":
		dst.MaxHeight = src.MaxHeight
	case "margin-top":
		dst.Margin[0] = src.Margin[0]
	case "margin-right":
		dst.Margin[1] = src.Margin[1]
	case "margin-bottom":
		dst.Margin[2] = src.Margin[2]
	case "margin-left":
		dst.Margin[3] = src.Margin[3]
	case "padding-top":
		dst.Padding[0] = src.Padding[0]
	case "padding-right":
		dst.Padding[1] = src.Padding[1]
	case "padding-bottom":
		dst.Padding[2] = src.Padding[2]
	case "padding-left":
		dst.Padding[3] = src.Padding[3]
	case "top":
		dst.Top = src.Top
	case "right":
		dst.Right = src.Right
	case "bottom":
		dst.Bottom = src.Bottom
	case "left":
		dst.Left = src.Left
	case "font-size":
		dst.FontSize = src.FontSize
	case "letter-spacing":
		dst.LetterSpacing = src.LetterSpacing
	case "word-spacing":
		dst.WordSpacing = src.WordSpacing
	case "line-height":
		dst.LineHeight = src.LineHeight
	}
}
