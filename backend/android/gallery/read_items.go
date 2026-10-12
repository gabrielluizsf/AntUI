//go:build android

package gallery

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// readItems walks the cursor, up to limit items when limit is above zero.
func readItems(e *jni.Env, cursor, collection jni.Object, modern bool, limit int) ([]Item, error) {
	var out []Item
	for {
		more, err := e.InvokeBool(cursor, "moveToNext", jni.Sig(jni.TBool))
		if err != nil || !more {
			return out, err
		}
		item, err := readItem(e, cursor, collection, modern)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
		if limit > 0 && len(out) >= limit {
			return out, nil
		}
	}
}
