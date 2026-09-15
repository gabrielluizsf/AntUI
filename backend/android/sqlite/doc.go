// Package sqlite is the SQLite every Android device already has.
//
// It goes through Java, which is the only way: the library is on the device
// and the NDK does not expose it, and linking a copy into the app would mean
// two SQLites open on one file, each with a cache the other knows nothing
// about — which is how a database gets corrupted.
//
// The API is database/sql's shape without being a driver: Open, Exec, Query,
// Next, Scan, Close. What it is not is a driver for database/sql, because
// the platform's own query API binds every argument as text and a driver
// that pretended otherwise would be lying to everything built on it.
package sqlite
