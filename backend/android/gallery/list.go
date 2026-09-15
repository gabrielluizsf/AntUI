//go:build android

package gallery

import (
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Item is one thing in the device's library.
type Item struct {
	// URI is its address, and is what [antui/backend/android/picker.Open] takes.
	URI string
	// Name is what the user sees.
	Name string
	// Size is in bytes, and Mime is what it holds.
	Size int64
	Mime string
	// Added is when it appeared in the library.
	Added time.Time
	// Duration is how long it plays, and is zero for a picture.
	Duration time.Duration
	// Album is the folder inside the library, on Android 10 and up.
	Album string
}

// Query is how to narrow a listing.
type Query struct {
	// Album keeps only what is in one folder — the same name [Create] took.
	// Empty means everything.
	Album string
	// Limit is how many at most, newest first. Zero means everything, which
	// on a real phone can be tens of thousands.
	Limit int
	// Mine keeps only what this app put there. It is the one listing that
	// needs no permission at all: an app can always see its own.
	Mine bool
}

// The columns read. They are the same names on every kind of media, because
// they are MediaStore.MediaColumns and not the per-kind ones.
var listColumns = []string{
	"_id", "_display_name", "_size", "mime_type", "date_added", "duration",
}

// List is what is in one of the device's libraries.
//
// **Seeing other apps' media needs a permission** — READ_EXTERNAL_STORAGE
// below Android 13, and READ_MEDIA_IMAGES, _VIDEO or _AUDIO from 13 — and
// without one the answer is only what this app put there, with no error to
// say so. [Query.Mine] asks for that on purpose, and is what an app that
// only wants its own should use: it needs nothing granted.
func List(kind Kind, q Query) ([]Item, error) {
	a := app.Current()
	if a == nil {
		return nil, app.ErrNoUIThread
	}
	modern := a.Config().SDK >= 29

	var out []Item
	err := jni.Do(func(e *jni.Env) error {
		columns := listColumns
		if modern {
			columns = append(append([]string{}, columns...), colPath)
		}
		projection, err := e.NewStrings(columns)
		if err != nil {
			return err
		}

		var where jni.Object
		var args jni.Object
		switch {
		case q.Mine:
			// Everything MediaStore reports to an app with no permission is
			// already its own, but saying so makes the intent plain and
			// keeps the answer the same when a permission is later added.
			if where, err = e.String("owner_package_name = ?"); err != nil {
				return err
			}
			name, err := e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
			if err != nil {
				return err
			}
			if args, err = e.NewStrings([]string{name}); err != nil {
				return err
			}
		case q.Album != "" && modern:
			if where, err = e.String(colPath + " LIKE ?"); err != nil {
				return err
			}
			if args, err = e.NewStrings([]string{"%/" + q.Album + "/%"}); err != nil {
				return err
			}
		}
		// Newest first, which is what a gallery shows and what a limit
		// should keep.
		order, err := e.String("date_added DESC")
		if err != nil {
			return err
		}

		class, field := kind.collection()
		collection, err := e.Constant(class, field, jni.TClass("android/net/Uri"))
		if err != nil {
			return err
		}
		resolver, err := e.Invoke(app.Context(), "getContentResolver",
			jni.Sig(jni.TClass("android/content/ContentResolver")))
		if err != nil {
			return err
		}
		cursor, err := e.Invoke(resolver, "query",
			jni.Sig(jni.TClass("android/database/Cursor"),
				jni.TClass("android/net/Uri"), jni.TArray(jni.TString),
				jni.TString, jni.TArray(jni.TString), jni.TString),
			jni.Ref(collection), jni.Ref(projection), jni.Ref(where),
			jni.Ref(args), jni.Ref(order))
		if err != nil || cursor.IsNil() {
			return err
		}
		defer e.InvokeVoid(cursor, "close", jni.Sig(jni.TVoid))

		for {
			more, err := e.InvokeBool(cursor, "moveToNext", jni.Sig(jni.TBool))
			if err != nil || !more {
				return err
			}
			item, err := readItem(e, cursor, collection, modern)
			if err != nil {
				return err
			}
			out = append(out, item)
			if q.Limit > 0 && len(out) >= q.Limit {
				return nil
			}
		}
	})
	return out, err
}

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
