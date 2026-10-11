package canvas

// Scaled is the same font at a multiple of its size, made once and kept.
func (f *Face) Scaled(n int) *Face {
	if n <= 1 {
		return f
	}
	if f == nil {
		return BuiltinScaled(n)
	}
	if f.file == nil {
		return builtinAt(f.scaleOf() * float64(n) * float64(FontHeight))
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.derived == nil {
		f.derived = map[int]*Face{}
	}
	if d, ok := f.derived[n]; ok {
		return d
	}
	d, err := ParseFace(f.file.Data, f.size*float64(n))
	if err != nil {
		return f
	}
	if f.next != nil {
		d.next = f.next.Scaled(n)
	}
	f.derived[n] = d
	return d
}
