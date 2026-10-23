//go:build linux || freebsd || openbsd || netbsd || dragonfly || darwin

package dbus

import "testing"

// TestNotifyBodyIsPlainText keeps text a program passes on — a chat message,
// a file name somebody else chose — from being read as markup by a server
// that understands it: a link or an <img> in it would be drawn as one, where
// the same text on Windows and macOS shows as written.
func TestNotifyBodyIsPlainText(t *testing.T) {
	calls := fakeBus(t, returnID)

	body := `<a href="https://evil.example">click</a> <img src="file:///etc/passwd"/> & done`
	if err := Notify("antui-demo", "Message", body); err != nil {
		t.Fatalf("Notify: %v", err)
	}
	call := <-calls
	got := readStrings(t, call.signature, call.body)
	want := `&lt;a href="https://evil.example"&gt;click&lt;/a&gt; &lt;img src="file:///etc/passwd"/&gt; &amp; done`
	if len(got) < 4 || got[3] != want {
		t.Fatalf("body = %q, want the markup escaped: %q", got, want)
	}
}

// TestNotifySurvivesBytesTheBusRefuses keeps a NUL or a broken UTF-8 sequence
// in a title from getting the whole message refused: the bus drops a
// connection that sends either, and the notification with it.
func TestNotifySurvivesBytesTheBusRefuses(t *testing.T) {
	calls := fakeBus(t, returnID)

	if err := Notify("antui-demo", "a\x00b\xffc", "ok"); err != nil {
		t.Fatalf("Notify: %v", err)
	}
	call := <-calls
	got := readStrings(t, call.signature, call.body)
	if len(got) < 3 || got[2] != "ab�c" {
		t.Fatalf("summary = %q, want the NUL dropped and the broken byte replaced", got)
	}
}
