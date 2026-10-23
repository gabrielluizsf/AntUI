package dbus

import "fmt"

// awaitReply reads until the answer to the call numbered serial arrives,
// stepping over everything else the bus sends meanwhile — the answer to
// Hello, the NameAcquired signal that follows it. An ERROR answer comes back
// as an error carrying its name and message.
func (c *conn) awaitReply(serial uint32) (message, error) {
	for {
		m, err := c.readMessage()
		if err != nil {
			return message{}, fmt.Errorf("dbus: waiting for an answer: %w", err)
		}
		if m.replySerial != serial {
			continue
		}
		switch m.typ {
		case typeMethodReturn:
			return m, nil
		case typeError:
			return message{}, fmt.Errorf("dbus: %s: %s", m.errorName, m.firstString())
		}
	}
}
