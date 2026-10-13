package sdk

import "runtime"

// exe adds the extension a program has on this machine.
func exe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
