package shim

import (
	"bytes"
	"io/fs"
	"testing"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// The dex in the repository has to be what the Java beside it compiles to.
//
// This is checkable because d8 is deterministic: the same class files give
// the same bytes every time. A committed artefact that cannot be checked
// against its source is a place for something to rot, and the whole reason
// the dex is committed is so that building an app needs no JDK — so the one
// machine that does have a JDK should be the one that notices.
func TestTheCommittedDexIsWhatTheSourceCompilesTo(t *testing.T) {
	tc, err := sdk.Find(sdk.Require{JDK: true})
	if err != nil {
		t.Skipf("no toolchain to rebuild the shim with: %v", err)
	}
	built, err := Build(tc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(built, Dex) {
		t.Errorf("classes.dex is %d bytes and the source compiles to %d; "+
			"rebuild it with \"antuiapk shim\"", len(Dex), len(built))
	}
}

func TestTheShimIsThere(t *testing.T) {
	if len(Dex) == 0 {
		t.Fatal("classes.dex is empty")
	}
	// Every dex file starts with this, and the version follows it.
	if !bytes.HasPrefix(Dex, []byte("dex\n")) {
		t.Errorf("classes.dex does not start like a dex file: %q", Dex[:8])
	}
	activity, err := Sources.ReadFile("AntuiActivity.java")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(activity, []byte("class AntuiActivity extends NativeActivity")) {
		t.Error("the Java source is not the shim")
	}
	names, err := fs.Glob(Sources, "*.java")
	if err != nil || len(names) < 2 {
		t.Errorf("the shim is %d files, and there should be one per listener kind", len(names))
	}
}
