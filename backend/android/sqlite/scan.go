//go:build android

package sqlite

import (
	"errors"
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Scan reads the current row into the values pointed at.
//
// It takes *string, *int64, *int, *float64, *bool, *[]byte and *any. A NULL
// leaves the destination at its zero value, and an *any receives nil — which
// is how a NULL is told from an empty string.
func (r *Rows) Scan(dest ...any) error {
	if r.closed {
		return errors.New("sqlite: the rows are closed")
	}
	if len(dest) > len(r.columns) && len(r.columns) > 0 {
		return fmt.Errorf("sqlite: %d destinations for %d columns",
			len(dest), len(r.columns))
	}
	return jni.Do(func(e *jni.Env) error {
		for i, d := range dest {
			if err := r.scanOne(e, int32(i), d); err != nil {
				return err
			}
		}
		return nil
	})
}
