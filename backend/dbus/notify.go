package dbus

import (
	"fmt"
	"time"
)

// notifyTimeout bounds the whole exchange. The bus and the notification
// server are both local, so an answer that has not come by then is a server
// that is not going to give one.
const notifyTimeout = 2 * time.Second

// Notify shows a desktop notification through the freedesktop notification
// server on the session bus, waiting for the server to take it. appName is
// who the notification says it is from.
//
// Every call is a connection of its own: notifications are rare, and a bus
// connection kept open for them would be one more thing to keep alive and
// reconnect across a session restart for nothing.
func Notify(appName, summary, body string) error {
	c, err := dial(time.Now().Add(notifyTimeout))
	if err != nil {
		return err
	}
	defer c.Close()

	// The bus wants Hello before anything else. Its answer need not be
	// waited for: the bus handles one connection's messages in order, so
	// Notify goes out right behind it.
	hello := methodCall(c.nextSerial(), "org.freedesktop.DBus", "/org/freedesktop/DBus",
		"org.freedesktop.DBus", "Hello", "", nil)
	serial := c.nextSerial()
	notify := methodCall(serial, "org.freedesktop.Notifications", "/org/freedesktop/Notifications",
		"org.freedesktop.Notifications", "Notify", notifySignature, notifyBody(appName, summary, body))
	if _, err := c.Write(append(hello, notify...)); err != nil {
		return fmt.Errorf("dbus: sending the notification: %w", err)
	}
	_, err = c.awaitReply(serial)
	return err
}
