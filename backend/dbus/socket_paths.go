package dbus

import (
	"net/url"
	"strings"
)

// socketPaths turns a bus address into the unix sockets it names, in the
// order to try them. An address is a list separated by ';', each entry a
// transport and its key=value options; only unix transports are kept, and an
// abstract socket comes back with the leading NUL Go spells it with.
func socketPaths(address string) []string {
	var paths []string
	for _, entry := range strings.Split(address, ";") {
		options, ok := strings.CutPrefix(entry, "unix:")
		if !ok {
			continue
		}
		for _, option := range strings.Split(options, ",") {
			key, value, _ := strings.Cut(option, "=")
			// Values are percent-escaped; an escape that does not decode is
			// a broken address, and its entry is skipped.
			value, err := url.PathUnescape(value)
			if err != nil {
				break
			}
			switch key {
			case "path":
				paths = append(paths, value)
			case "abstract":
				paths = append(paths, "\x00"+value)
			}
		}
	}
	return paths
}
