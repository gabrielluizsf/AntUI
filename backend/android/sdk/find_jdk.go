package sdk

import (
	"os"
	"os/exec"
	"path/filepath"
)

func findJDK() (JDK, bool) {
	if home := os.Getenv("JAVA_HOME"); home != "" {
		javac := filepath.Join(home, "bin", exe("javac"))
		if _, err := os.Stat(javac); err == nil {
			return JDK{Dir: home, Javac: javac}, true
		}
	}
	javac, err := exec.LookPath("javac")
	if err != nil {
		return JDK{}, false
	}
	// Follow the symlink an alternatives system leaves in /usr/bin, so the
	// home we report is the real one and not /usr.
	if real, err := filepath.EvalSymlinks(javac); err == nil {
		javac = real
	}
	return JDK{Dir: filepath.Dir(filepath.Dir(javac)), Javac: javac}, true
}
