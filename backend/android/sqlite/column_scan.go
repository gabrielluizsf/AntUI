//go:build android

package sqlite

import "fmt"

// scan reads the column into dest, a pointer to one of the types a Go caller
// uses. A NULL becomes the zero of that type.
func (c column) scan(kind int32, dest any) error {
	switch d := dest.(type) {
	case *string:
		if kind == typeNull {
			*d = ""
			return nil
		}
		v, err := c.text()
		*d = v
		return err
	case *int64:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := c.whole()
		*d = v
		return err
	case *int:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := c.whole()
		*d = int(v)
		return err
	case *float64:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := c.real()
		*d = v
		return err
	case *bool:
		return c.scanBool(kind, d)
	case *[]byte:
		if kind == typeNull {
			*d = nil
			return nil
		}
		v, err := c.blob()
		*d = v
		return err
	case *any:
		return c.scanAny(kind, d)
	}
	return fmt.Errorf("sqlite: cannot scan into %T", dest)
}
