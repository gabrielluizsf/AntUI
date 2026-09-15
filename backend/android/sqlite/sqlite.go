//go:build android

package sqlite

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// modePrivate is Context.MODE_PRIVATE, and the only mode still allowed.
const modePrivate = 0

// DB is an open database.
//
// SQLite is on every Android device and has been since the first one, but
// only through Java: the C library is there and the NDK does not expose it,
// and linking a copy of SQLite into the app would mean shipping a second one
// beside the platform's and having two caches of the same file.
type DB struct {
	obj    jni.Object
	closed bool
}

// Open opens or creates a database in the app's own directory, by name.
// It is the ordinary case: the file lives where the system backs it up and
// deletes it with the app.
func Open(name string) (*DB, error) {
	if name == "" {
		return nil, errors.New("sqlite: a database needs a name")
	}
	var db DB
	err := jni.Do(func(e *jni.Env) error {
		jname, err := e.String(name)
		if err != nil {
			return err
		}
		obj, err := e.Invoke(app.Context(), "openOrCreateDatabase",
			jni.Sig(jni.TClass("android/database/sqlite/SQLiteDatabase"),
				jni.TString, jni.TInt,
				jni.TClass("android/database/sqlite/SQLiteDatabase$CursorFactory")),
			jni.Ref(jname), jni.Int(modePrivate), jni.Ref(jni.Object{}))
		if err != nil {
			return err
		}
		if obj.IsNil() {
			return fmt.Errorf("sqlite: cannot open %s", name)
		}
		db.obj = e.Global(obj)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// Exec runs a statement that answers nothing: a CREATE, an INSERT, an
// UPDATE. The arguments replace the ? marks, in order.
func (db *DB) Exec(sql string, args ...any) error {
	if db.closed {
		return errors.New("sqlite: the database is closed")
	}
	return jni.Do(func(e *jni.Env) error {
		jsql, err := e.String(sql)
		if err != nil {
			return err
		}
		if len(args) == 0 {
			return e.InvokeVoid(db.obj, "execSQL",
				jni.Sig(jni.TVoid, jni.TString), jni.Ref(jsql))
		}
		boxed, err := boxAll(e, args)
		if err != nil {
			return err
		}
		return e.InvokeVoid(db.obj, "execSQL",
			jni.Sig(jni.TVoid, jni.TString, jni.TArray(jni.TObject)),
			jni.Ref(jsql), jni.Ref(boxed))
	})
}

// Query runs a SELECT.
//
// **Every argument is bound as text.** That is the platform's API and not a
// choice here: rawQuery takes a String array and nothing else. SQLite's own
// typing means "42" compares equal to 42 in almost every case, but a column
// declared with no affinity will hold the text — so a number that must stay
// a number belongs in the statement or in an Exec.
func (db *DB) Query(sql string, args ...any) (*Rows, error) {
	if db.closed {
		return nil, errors.New("sqlite: the database is closed")
	}
	rows := &Rows{}
	err := jni.Do(func(e *jni.Env) error {
		jsql, err := e.String(sql)
		if err != nil {
			return err
		}
		var jargs jni.Object
		if len(args) > 0 {
			text := make([]string, len(args))
			for i, a := range args {
				text[i] = asText(a)
			}
			if jargs, err = e.NewStrings(text); err != nil {
				return err
			}
		}
		cursor, err := e.Invoke(db.obj, "rawQuery",
			jni.Sig(jni.TClass("android/database/Cursor"), jni.TString,
				jni.TArray(jni.TString)),
			jni.Ref(jsql), jni.Ref(jargs))
		if err != nil {
			return err
		}
		if cursor.IsNil() {
			return errors.New("sqlite: the query answered nothing at all")
		}
		rows.cursor = e.Global(cursor)
		names, err := e.Invoke(rows.cursor, "getColumnNames",
			jni.Sig(jni.TArray(jni.TString)))
		if err == nil && !names.IsNil() {
			rows.columns, _ = e.GoStrings(names)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Transaction runs f inside one, committing when it returns nil and undoing
// everything when it does not.
//
// It is worth using for anything that writes more than once: SQLite commits
// each statement on its own otherwise, and each commit waits for the disk.
// A thousand inserts inside a transaction take about as long as one.
//
// **The goroutine is pinned to its thread for the whole of it**, and that is
// not a detail. A SQLiteDatabase transaction belongs to the thread that
// began it, and Go moves a goroutine between threads whenever it likes — so
// without this the begin happens on one thread and the first statement on
// another, which then waits thirty seconds for a lock the first thread is
// holding and gives up. It does not fail: it hangs, and the log says
// something about a connection pool.
func (db *DB) Transaction(f func() error) error {
	release, err := jni.Hold()
	if err != nil {
		return err
	}
	defer release()

	if err := db.call("beginTransaction"); err != nil {
		return err
	}
	err = f()
	if err == nil {
		if markErr := db.call("setTransactionSuccessful"); markErr != nil {
			err = markErr
		}
	}
	if endErr := db.call("endTransaction"); endErr != nil && err == nil {
		err = endErr
	}
	return err
}

func (db *DB) call(method string) error {
	return jni.Do(func(e *jni.Env) error {
		return e.InvokeVoid(db.obj, method, jni.Sig(jni.TVoid))
	})
}

// Path is where the file is, for a caller that wants to copy or back it up.
func (db *DB) Path() (string, error) {
	var out string
	err := jni.Do(func(e *jni.Env) error {
		var err error
		out, err = e.InvokeString(db.obj, "getPath", jni.Sig(jni.TString))
		return err
	})
	return out, err
}

// Close shuts the database. Rows opened from it must be closed first.
func (db *DB) Close() error {
	if db.closed {
		return nil
	}
	db.closed = true
	return jni.Do(func(e *jni.Env) error {
		e.InvokeVoid(db.obj, "close", jni.Sig(jni.TVoid))
		e.DeleteGlobal(db.obj)
		return nil
	})
}

// Rows is the answer to a query, walked one row at a time.
type Rows struct {
	cursor  jni.Object
	columns []string
	err     error
	closed  bool
}

// Columns is what the query answered with.
func (r *Rows) Columns() []string { return r.columns }

// Next moves to the next row and reports whether there was one.
func (r *Rows) Next() bool {
	if r.closed || r.err != nil {
		return false
	}
	var more bool
	if err := jni.Do(func(e *jni.Env) error {
		var err error
		more, err = e.InvokeBool(r.cursor, "moveToNext", jni.Sig(jni.TBool))
		return err
	}); err != nil {
		r.err = err
		return false
	}
	return more
}

// The types a column can hold, as Cursor numbers them.
const (
	typeNull = 0
	typeInt  = 1
	typeReal = 2
	typeText = 3
	typeBlob = 4
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

func (r *Rows) scanOne(e *jni.Env, i int32, dest any) error {
	kind, err := e.InvokeInt(r.cursor, "getType", jni.Sig(jni.TInt, jni.TInt), jni.Int(i))
	if err != nil {
		return err
	}
	text := func() (string, error) {
		return e.InvokeString(r.cursor, "getString", jni.Sig(jni.TString, jni.TInt), jni.Int(i))
	}
	whole := func() (int64, error) {
		return e.InvokeLong(r.cursor, "getLong", jni.Sig(jni.TLong, jni.TInt), jni.Int(i))
	}
	real := func() (float64, error) {
		v, err := e.InvokeFloat(r.cursor, "getFloat", jni.Sig(jni.TFloat, jni.TInt), jni.Int(i))
		return float64(v), err
	}
	blob := func() ([]byte, error) {
		o, err := e.Invoke(r.cursor, "getBlob",
			jni.Sig(jni.TArray(jni.TByte), jni.TInt), jni.Int(i))
		if err != nil || o.IsNil() {
			return nil, err
		}
		return e.GoBytes(o)
	}

	switch d := dest.(type) {
	case *string:
		if kind == typeNull {
			*d = ""
			return nil
		}
		v, err := text()
		*d = v
		return err
	case *int64:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := whole()
		*d = v
		return err
	case *int:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := whole()
		*d = int(v)
		return err
	case *float64:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := real()
		*d = v
		return err
	case *bool:
		if kind == typeNull {
			*d = false
			return nil
		}
		v, err := whole()
		*d = v != 0
		return err
	case *[]byte:
		if kind == typeNull {
			*d = nil
			return nil
		}
		v, err := blob()
		*d = v
		return err
	case *any:
		switch kind {
		case typeNull:
			*d = nil
		case typeInt:
			v, err := whole()
			*d = v
			return err
		case typeReal:
			v, err := real()
			*d = v
			return err
		case typeBlob:
			v, err := blob()
			*d = v
			return err
		default:
			v, err := text()
			*d = v
			return err
		}
		return nil
	}
	return fmt.Errorf("sqlite: cannot scan into %T", dest)
}

// Err is what stopped the walk, or nil.
func (r *Rows) Err() error { return r.err }

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

// boxAll turns Go values into the Objects execSQL binds.
func boxAll(e *jni.Env, args []any) (jni.Object, error) {
	cls, err := e.Class("java/lang/Object")
	if err != nil {
		return jni.Object{}, err
	}
	arr, err := e.NewObjects(cls, len(args))
	if err != nil {
		return jni.Object{}, err
	}
	err = e.Frame(len(args)*4+16, func() error {
		for i, a := range args {
			boxed, err := box(e, a)
			if err != nil {
				return err
			}
			if err := e.SetIndex(arr, i, boxed); err != nil {
				return err
			}
		}
		return nil
	})
	return arr, err
}

func box(e *jni.Env, v any) (jni.Object, error) {
	switch t := v.(type) {
	case nil:
		return jni.Object{}, nil
	case string:
		return e.String(t)
	case []byte:
		return e.NewBytes(t)
	case bool:
		n := int64(0)
		if t {
			n = 1
		}
		return e.Static("java/lang/Long", "valueOf",
			jni.Sig(jni.TClass("java/lang/Long"), jni.TLong), jni.Long(n))
	case int:
		return box(e, int64(t))
	case int32:
		return box(e, int64(t))
	case int64:
		return e.Static("java/lang/Long", "valueOf",
			jni.Sig(jni.TClass("java/lang/Long"), jni.TLong), jni.Long(t))
	case float32:
		return box(e, float64(t))
	case float64:
		return e.Static("java/lang/Double", "valueOf",
			jni.Sig(jni.TClass("java/lang/Double"), jni.TDouble), jni.Double(t))
	}
	return jni.Object{}, fmt.Errorf("sqlite: cannot bind %T", v)
}

// asText is how a query argument is written, since the platform binds them
// all as strings.
func asText(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "1"
		}
		return "0"
	case int:
		return strconv.Itoa(t)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case []byte:
		return string(t)
	}
	return fmt.Sprint(v)
}
