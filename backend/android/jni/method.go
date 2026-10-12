//go:build android

package jni

// method finds a method on an object's own class.
//
// **It does not cache, and it lets the class reference go before it
// returns.** Both halves of that were wrong before and both were invisible.
//
// [Env.ClassOf] hands back a *local* reference, and a local reference is a
// slot in a table that is reused as soon as the frame that made it ends — so
// the same number means one class now and a different class a minute later.
// Keying a cache on it is how `close()V`, resolved once for an
// SQLiteCursor, came back for an SQLiteDatabase. Nothing notices: the ids
// are opaque, the call goes through, and the VM does whatever that method's
// code does to an object of the wrong type. The only thing that says so is
// the VM's own CheckJNI —
//
//	adb shell setprop debug.checkjni 1 && adb shell stop && adb shell start
//
// which aborts with "can't call void android.database.sqlite.
// SQLiteCursor.close() on instance of android.database.sqlite.
// SQLiteDatabase". It is worth turning on now and then for exactly this.
//
// Not deleting the reference was the same mistake's other half: a loop of
// Invoke calls filled the local reference table with class references nobody
// would look at again, and the table is 512 entries deep.
//
// The named-class calls — [Env.Static], [Env.Constant], [Env.Make] — do
// cache, and safely: [Env.Class] keeps one **global** reference per name,
// and a global is unique and permanent. A caller who wants a cached lookup
// on a hot path should name the class.
func (e *Env) method(o Object, name, sig string) (Method, error) {
	c, err := e.ClassOf(o)
	if err != nil {
		return Method{}, err
	}
	defer e.Delete(c.Object())
	return e.Method(c, name, sig)
}
