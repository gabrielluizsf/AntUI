package svg

// rootStyle is where the style of everything starts for this drawing: [baseStyle]
// with the address the drawing was read from, which a relative address inside it
// resolves against. A drawing read from no file has none, and then an address is
// left as it was written. Every picture an element of the drawing paints — the
// shape itself, what a `<pattern>`, a `<mask>`, a `<clipPath>` or a `<marker>`
// holds — starts from here, so all of them resolve against the same base.
func (img *Image) rootStyle() Style {
	st := baseStyle()
	st.base = img.base
	return st
}
