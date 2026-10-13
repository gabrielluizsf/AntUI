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
		dst.Margin[SideTop] = src.Margin[SideTop]
	case "margin-right":
		dst.Margin[SideRight] = src.Margin[SideRight]
	case "margin-bottom":
		dst.Margin[SideBottom] = src.Margin[SideBottom]
	case "margin-left":
		dst.Margin[SideLeft] = src.Margin[SideLeft]
	case "padding-top":
		dst.Padding[SideTop] = src.Padding[SideTop]
	case "padding-right":
		dst.Padding[SideRight] = src.Padding[SideRight]
	case "padding-bottom":
		dst.Padding[SideBottom] = src.Padding[SideBottom]
	case "padding-left":
		dst.Padding[SideLeft] = src.Padding[SideLeft]
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
