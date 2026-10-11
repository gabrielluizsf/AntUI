package canvas

func (c *Color) UnmarshalText(text []byte) error {
	v, err := ParseColor(string(text))
	if err != nil {
		return err
	}
	*c = v
	return nil
}
