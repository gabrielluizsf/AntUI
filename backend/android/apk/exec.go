package apk

import (
	"fmt"
	"os/exec"
	"strings"
)

// run executes a tool and returns what it printed. On failure the error
// carries the output, because every one of these tools says what is wrong
// there and nowhere else — an exit status on its own is never enough to act
// on.
func run(name string, args ...string) (string, error) {
	return runIn("", nil, name, args...)
}

// runIn is run with a working directory and extra environment.
func runIn(dir string, env []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = env
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s %s: %w\n%s",
			shortName(name), strings.Join(args, " "), err, out)
	}
	return string(out), nil
}

// shortName is the tool's own name, so an error reads "aapt2 link: ..." and
// not the whole path to it.
func shortName(path string) string {
	if i := strings.LastIndexAny(path, `/\`); i >= 0 {
		return path[i+1:]
	}
	return path
}
