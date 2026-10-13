//go:build android

package gallery

import "github.com/gabrielluizsf/antui/backend/android/jni"

// buildValues is the ContentValues a MediaStore insert takes.
func buildValues(e *jni.Env, kind Kind, name, mime, album string, modern bool) (jni.Object, error) {
	values, err := e.Make("android/content/ContentValues", jni.Sig(jni.TVoid))
	if err != nil {
		return jni.Object{}, err
	}
	put := func(key, value string) error { return putString(e, values, key, value) }
	if err := put(colName, name); err != nil {
		return jni.Object{}, err
	}
	if err := put(colMime, mime); err != nil {
		return jni.Object{}, err
	}
	if modern {
		path := kind.folder()
		if album != "" {
			path += "/" + album
		}
		if err := put(colPath, path); err != nil {
			return jni.Object{}, err
		}
		// Pending until it is written, so that nothing sees a half file.
		if err := putInt(e, values, colPending, 1); err != nil {
			return jni.Object{}, err
		}
	}
	return values, nil
}
