//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Close gives the cursor back. **It has to be called**: a cursor left open
// holds a window of the database's memory, and the platform complains loudly
// in the log when one is collected without being closed.
func (r *Rows) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	return jni.Do(func(e *jni.Env) error {
		e.InvokeVoid(r.cursor, "close", jni.Sig(jni.TVoid))
		e.DeleteGlobal(r.cursor)
		return nil
	})
}
