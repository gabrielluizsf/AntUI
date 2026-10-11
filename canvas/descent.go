package canvas

func (f *Face) Descent() int {
	if f == nil {
		return builtin.descent
	}
	return f.descent
}
