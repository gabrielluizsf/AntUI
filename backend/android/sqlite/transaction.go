//go:build android

package sqlite

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
