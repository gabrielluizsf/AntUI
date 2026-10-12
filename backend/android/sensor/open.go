//go:build android

package sensor

import (
	"errors"
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// Open starts the named sensors.
func Open(kinds []Kind, opt Options) (*Stream, error) {
	if len(kinds) == 0 {
		return nil, errors.New("sensor: no sensors asked for")
	}
	m, err := sensors()
	if err != nil {
		return nil, err
	}
	chosen := make([]ndk.Sensor, 0, len(kinds))
	for _, k := range kinds {
		s := m.Default(int(k))
		if opt.Wake {
			s = m.DefaultWake(int(k), true)
		}
		if s.IsNil() {
			return nil, fmt.Errorf("sensor: this device has no %s", k)
		}
		chosen = append(chosen, s)
	}

	st := &Stream{
		latest:  make(map[Kind]Reading, len(kinds)),
		fresh:   make(map[Kind]bool, len(kinds)),
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	if opt.Buffer > 0 {
		st.events = make(chan Reading, opt.Buffer)
	}

	ready := make(chan error, 1)
	go st.run(m, chosen, opt, ready)
	if err := <-ready; err != nil {
		return nil, err
	}
	return st, nil
}
