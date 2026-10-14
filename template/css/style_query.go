package css

// StyleQuery is the set of inputs that decides one element's cascade style.
type StyleQuery struct {
	Tag     string
	Classes []string
	State   State
	Viewport
	Font int // the element's font size (for em/ch/ex)
	Root int // the document font size (for rem)
}

// Units returns the measurement context as a Units value.
func (q StyleQuery) Units() Units {
	return Units{
		Width:  q.Viewport.Width,
		Height: q.Viewport.Height,
		Font:   q.Font,
		Root:   q.Root,
		Scale:  q.Scale,
		Scheme: q.Scheme,
	}
}

// PropertyQuery is the set of inputs for reading a single property value.
type PropertyQuery struct {
	Tag     string
	Classes []string
	Prop    string
	Viewport
	Font int
	Root int
}

// Units returns the measurement context as a Units value.
func (q PropertyQuery) Units() Units {
	return Units{
		Width:  q.Viewport.Width,
		Height: q.Viewport.Height,
		Font:   q.Font,
		Root:   q.Root,
		Scale:  q.Scale,
		Scheme: q.Scheme,
	}
}
