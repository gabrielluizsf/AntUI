package shim

import (
	"fmt"
	"os/exec"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// compileJava compiles the sources against android.jar as the bootclasspath
// rather than the JDK's own: the shim extends a platform class that no JDK
// has.
func compileJava(tc *sdk.Toolchain, classes string, srcs []string) error {
	args := []string{
		"-source", SourceLevel, "-target", SourceLevel,
		"-bootclasspath", tc.Platform.Jar(),
		"-nowarn",
		"-d", classes,
	}
	out, err := exec.Command(tc.JDK.Javac, append(args, srcs...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("shim: javac: %w\n%s", err, out)
	}
	return nil
}
