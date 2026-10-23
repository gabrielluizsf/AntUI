//go:build linux || freebsd || openbsd || netbsd || dragonfly || darwin

package dbus

import (
	"bufio"
	"encoding/binary"
	"net"
	"path/filepath"
	"strings"
	"testing"
)

// fakeBus is a session bus with one notification server on it: it
// authenticates whoever connects, answers Hello, and answers Notify with
// what reply builds from the call it got.
func fakeBus(t *testing.T, reply func(call message) []byte) <-chan message {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bus")
	ln, err := net.Listen("unix", path)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+path+",guid=0123")

	calls := make(chan message, 1)
	go func() {
		raw, err := ln.Accept()
		if err != nil {
			return
		}
		defer raw.Close()
		c := &conn{Conn: raw, r: bufio.NewReader(raw)}
		if b, _ := c.r.ReadByte(); b != 0 {
			return
		}
		if line, _ := c.r.ReadString('\n'); !strings.HasPrefix(line, "AUTH EXTERNAL ") {
			return
		}
		raw.Write([]byte("OK 0123\r\n"))
		if line, _ := c.r.ReadString('\n'); line != "BEGIN\r\n" {
			return
		}
		for {
			m, err := c.readMessage()
			if err != nil {
				return
			}
			switch m.member {
			case "Hello":
				var body encoder
				body.string(":1.42")
				raw.Write(encodeMessage(typeMethodReturn, 1, []headerField{
					{code: fieldReplySerial, sig: 'u', num: m.serial},
				}, "s", body.buf))
			case "Notify":
				calls <- m
				raw.Write(reply(m))
			}
		}
	}()
	return calls
}

// returnID is a notification server taking the call and giving it id 7.
func returnID(call message) []byte {
	var body encoder
	body.uint32(7)
	return encodeMessage(typeMethodReturn, 2, []headerField{
		{code: fieldReplySerial, sig: 'u', num: call.serial},
	}, "u", body.buf)
}

// readStrings decodes the leading run of strings and uint32s in a body laid
// out as sig says, which is all of Notify's arguments that carry anything.
func readStrings(t *testing.T, sig string, body []byte) []string {
	t.Helper()
	var out []string
	i := 0
	for _, c := range sig {
		i = (i + 3) &^ 3
		switch c {
		case 's':
			n := int(binary.LittleEndian.Uint32(body[i:]))
			out = append(out, string(body[i+4:i+4+n]))
			i += 4 + n + 1
		case 'u':
			i += 4
		default:
			return out
		}
	}
	return out
}

func TestNotifySendsTitleAndBody(t *testing.T) {
	calls := fakeBus(t, returnID)

	if err := Notify("antui-demo", "Download", "report.pdf está pronto"); err != nil {
		t.Fatalf("Notify: %v", err)
	}
	call := <-calls
	if call.signature != notifySignature {
		t.Fatalf("signature = %q, want %q", call.signature, notifySignature)
	}
	got := readStrings(t, call.signature, call.body)
	want := []string{"antui-demo", "", "Download", "report.pdf está pronto"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("arguments = %q, want %q", got, want)
	}
	// The body ends with the timeout: -1, the server's own.
	if tail := int32(binary.LittleEndian.Uint32(call.body[len(call.body)-4:])); tail != -1 {
		t.Fatalf("expire_timeout = %d, want -1", tail)
	}
}

func TestNotifyReportsTheServersError(t *testing.T) {
	fakeBus(t, func(call message) []byte {
		var body encoder
		body.string("no notification server")
		return encodeMessage(typeError, 2, []headerField{
			{code: fieldReplySerial, sig: 'u', num: call.serial},
			{code: fieldErrorName, sig: 's', str: "org.freedesktop.DBus.Error.ServiceUnknown"},
		}, "s", body.buf)
	})

	err := Notify("antui-demo", "t", "b")
	if err == nil || !strings.Contains(err.Error(), "ServiceUnknown") ||
		!strings.Contains(err.Error(), "no notification server") {
		t.Fatalf("err = %v, want the server's error name and message", err)
	}
}

func TestNotifyWithNoBus(t *testing.T) {
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+filepath.Join(t.TempDir(), "missing"))
	if err := Notify("antui-demo", "t", "b"); err == nil {
		t.Fatal("Notify with nothing listening should fail")
	}
}
