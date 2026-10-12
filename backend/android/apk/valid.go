package apk

import (
	"fmt"
)

// valid reports what is wrong with a link, if anything.
func (l DeepLink) valid() error {
	if l.Scheme == "" {
		return fmt.Errorf("apk: a link needs a scheme")
	}
	if l.Host == "" && (l.Scheme == "http" || l.Scheme == "https") {
		return fmt.Errorf("apk: a %s link needs a host, or it claims the whole web", l.Scheme)
	}
	if l.Verify && l.Host == "" {
		return fmt.Errorf("apk: only a link with a host can be verified")
	}
	return nil
}
