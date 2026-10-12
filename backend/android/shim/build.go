package shim

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Build compiles the shim from [Source] and returns the dex.
//
// It needs a JDK and the SDK's build tools, which building an app does not —
// which is the whole reason the result is committed. Nothing is written
// outside the temporary directory it makes.
func Build(tc *sdk.Toolchain) ([]byte, error) {
	if tc.JDK.Javac == "" {
		return nil, fmt.Errorf("shim: rebuilding needs a JDK, and none was found")
	}
	if tc.Platform.Dir == "" {
		return nil, fmt.Errorf("shim: rebuilding needs an SDK platform to compile against")
	}
	if tc.BuildTools.Dir == "" {
		return nil, fmt.Errorf("shim: rebuilding needs the SDK build tools for d8")
	}

	work, err := os.MkdirTemp("", "antui-shim-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(work)

	names, err := fs.Glob(Sources, "*.java")
	if err != nil || len(names) == 0 {
		return nil, fmt.Errorf("shim: no Java to compile")
	}
	sort.Strings(names) // so the class files reach d8 in a fixed order
	srcs := make([]string, 0, len(names))
	for _, name := range names {
		b, err := Sources.ReadFile(name)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(work, name)
		if err := os.WriteFile(path, b, 0o644); err != nil {
			return nil, err
		}
		srcs = append(srcs, path)
	}
	classes := filepath.Join(work, "classes")
	if err := os.MkdirAll(classes, 0o755); err != nil {
		return nil, err
	}

	// Compiled against android.jar as the bootclasspath rather than the
	// JDK's own: the shim extends a platform class that no JDK has.
	args := []string{
		"-source", SourceLevel, "-target", SourceLevel,
		"-bootclasspath", tc.Platform.Jar(),
		"-nowarn",
		"-d", classes,
	}
	out, err := exec.Command(tc.JDK.Javac, append(args, srcs...)...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("shim: javac: %w\n%s", err, out)
	}

	var compiled []string
	err = filepath.WalkDir(classes, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(p, ".class") {
			compiled = append(compiled, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	if len(compiled) == 0 {
		return nil, fmt.Errorf("shim: javac produced no class files")
	}

	// d8 is deterministic given the same inputs in the same order, which is
	// what makes the committed dex checkable.
	sort.Strings(compiled)
	d8 := append([]string{"--lib", tc.Platform.Jar(), "--output", work}, compiled...)
	if out, err := exec.Command(tc.BuildTools.D8(), d8...).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("shim: d8: %w\n%s", err, out)
	}
	return os.ReadFile(filepath.Join(work, "classes.dex"))
}
