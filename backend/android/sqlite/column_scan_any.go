//go:build android

package sqlite

// scanAny reads into an *any, where the type is whatever the column holds.
func (c column) scanAny(kind int32, d *any) error {
	switch kind {
	case typeNull:
		*d = nil
	case typeInt:
		v, err := c.whole()
		*d = v
		return err
	case typeReal:
		v, err := c.real()
		*d = v
		return err
	case typeBlob:
		v, err := c.blob()
		*d = v
		return err
	default:
		v, err := c.text()
		*d = v
		return err
	}
	return nil
}
