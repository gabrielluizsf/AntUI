//go:build android

package sensor

import (
	"time"
)

// List is every sensor the device has.
func List() ([]Info, error) {
	m, err := sensors()
	if err != nil {
		return nil, err
	}
	all := m.List()
	out := make([]Info, 0, len(all))
	for _, s := range all {
		out = append(out, Info{
			Kind:       Kind(s.Type()),
			Name:       s.Name(),
			Vendor:     s.Vendor(),
			Resolution: s.Resolution(),
			MinDelay:   time.Duration(s.MinDelay()) * time.Microsecond,
		})
	}
	return out, nil
}
