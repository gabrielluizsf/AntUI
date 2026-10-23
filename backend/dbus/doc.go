// Package dbus is just enough of the D-Bus wire protocol to put a
// notification on a freedesktop desktop: connect to the session bus,
// authenticate, make one method call and wait for its answer.
//
// It is not a D-Bus library. There is no signal handling, no object export
// and no general marshaller — only the handful of types the calls here send
// and the header fields their answers come back with. Like the X11 backend
// beside it, it is spoken over the socket directly so the library stays pure
// Go with no dependencies.
package dbus
