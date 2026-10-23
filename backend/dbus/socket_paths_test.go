package dbus

import (
	"slices"
	"testing"
)

func TestSocketPaths(t *testing.T) {
	cases := []struct {
		address string
		want    []string
	}{
		{"unix:path=/run/user/1000/bus", []string{"/run/user/1000/bus"}},
		{"unix:abstract=/tmp/dbus-Xy,guid=ab12", []string{"\x00/tmp/dbus-Xy"}},
		{"unix:path=/tmp/a%20b", []string{"/tmp/a b"}},
		{"tcp:host=localhost,port=1;unix:path=/x", []string{"/x"}},
		{"tcp:host=localhost,port=1", nil},
		{"", nil},
	}
	for _, c := range cases {
		if got := socketPaths(c.address); !slices.Equal(got, c.want) {
			t.Errorf("socketPaths(%q) = %q, want %q", c.address, got, c.want)
		}
	}
}

func TestBusAddressFallsBackToTheRuntimeDir(t *testing.T) {
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "")
	t.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	got, err := busAddress()
	if err != nil || got != "unix:path=/run/user/1000/bus" {
		t.Fatalf("busAddress() = %q, %v; want the runtime directory's bus", got, err)
	}
}

// TestNotifyBodyLayout pins the bytes of an empty Notify body, where the
// alignment rules are easiest to get wrong: the empty a{sv} still pads to 8
// after its length.
func TestNotifyBodyLayout(t *testing.T) {
	got := notifyBody("", "", "")
	// Offsets: app_name 0-4, pad, replaces_id 8-11, app_icon 12-16, pad,
	// summary 20-24, pad, body 28-32, pad, actions' length 36-39, hints'
	// length 40-43 then padding to 48, expire_timeout 48-51.
	want := make([]byte, 52)
	copy(want[48:], []byte{0xff, 0xff, 0xff, 0xff})
	if !slices.Equal(got, want) {
		t.Fatalf("notifyBody = % x\nwant          % x", got, want)
	}
}
