//go:build android

package gallery

import (
	"time"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func readItem(e *jni.Env, cursor, collection jni.Object, modern bool) (Item, error) {
	var it Item
	str := func(i int32) string {
		v, _ := e.InvokeString(cursor, "getString", jni.Sig(jni.TString, jni.TInt), jni.Int(i))
		return v
	}
	num := func(i int32) int64 {
		v, _ := e.InvokeLong(cursor, "getLong", jni.Sig(jni.TLong, jni.TInt), jni.Int(i))
		return v
	}
	id := num(0)
	it.Name = str(1)
	it.Size = num(2)
	it.Mime = str(3)
	// date_added is in seconds, not milliseconds, which is the one column in
	// MediaStore that is — everything else that looks like a time is in
	// milliseconds, and mixing them puts a photograph in 1970.
	if secs := num(4); secs > 0 {
		it.Added = time.Unix(secs, 0)
	}
	if ms := num(5); ms > 0 {
		it.Duration = time.Duration(ms) * time.Millisecond
	}
	if modern {
		it.Album = str(6)
	}

	// The address of the row, which is the collection with the id on the
	// end. There is no column that holds it.
	uri, err := e.Static("android/content/ContentUris", "withAppendedId",
		jni.Sig(jni.TClass("android/net/Uri"), jni.TClass("android/net/Uri"), jni.TLong),
		jni.Ref(collection), jni.Long(id))
	if err != nil {
		return it, err
	}
	it.URI, err = e.InvokeString(uri, "toString", jni.Sig(jni.TString))
	return it, err
}
