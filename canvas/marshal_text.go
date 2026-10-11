package canvas

// MarshalText and UnmarshalText let a colour be written down — in JSON, in a
// config file, in anything that goes through encoding.TextMarshaler — as
// "#3E63DD" rather than as the number 4283372509.
func (c Color) MarshalText() ([]byte, error) { return []byte(c.String()), nil }
