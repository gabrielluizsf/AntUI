package antui

import (
	"os"
	"path/filepath"
)

// fontDirs is where a font file lives when the registry names only the file:
// the system's own directory, and the per-user one Windows 10 added for
// fonts installed without administrator rights.
func fontDirs() []string {
	dirs := []string{`C:\Windows\Fonts`}
	if root := os.Getenv("SystemRoot"); root != "" {
		dirs[0] = filepath.Join(root, "Fonts")
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		dirs = append(dirs, filepath.Join(local, "Microsoft", "Windows", "Fonts"))
	}
	return dirs
}
