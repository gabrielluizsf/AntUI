//go:build android

package sqlite

// scanBool reads a whole number as a boolean: zero is false.
func (c column) scanBool(kind int32, d *bool) error {
	if kind == typeNull {
		*d = false
		return nil
	}
	v, err := c.whole()
	*d = v != 0
	return err
}
