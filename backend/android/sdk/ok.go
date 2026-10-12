package sdk

// OK reports whether anything found would stop a build.
func (d *Doctor) OK() bool {
	for _, f := range d.Findings {
		if f.Level == Bad {
			return false
		}
	}
	return true
}
