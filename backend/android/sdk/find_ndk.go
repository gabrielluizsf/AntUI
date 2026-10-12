package sdk

import (
	"os"
	"path/filepath"
)

func findNDK(root string) (NDK, bool) {
	var found []NDK
	add := func(dir string) {
		if dir == "" {
			return
		}
		n, ok := readNDK(dir)
		if ok {
			found = append(found, n)
		}
	}
	for _, env := range []string{"ANDROID_NDK_HOME", "ANDROID_NDK_ROOT", "NDK_ROOT"} {
		add(os.Getenv(env))
	}
	if len(found) > 0 {
		return found[0], true // an environment that names one means it
	}
	for _, name := range subdirs(filepath.Join(root, "ndk")) {
		add(filepath.Join(root, "ndk", name))
	}
	add(filepath.Join(root, "ndk-bundle"))
	return pick(found, func(n NDK) Version { return n.Version })
}
