//go:build android

package camera

import "github.com/gabrielluizsf/antui/backend/android/ndk"

// List is every camera the app may open.
func List() ([]Info, error) {
	m, err := ndk.Cameras()
	if err != nil {
		return nil, err
	}
	defer m.Close()
	ids, err := m.IDs()
	if err != nil {
		return nil, err
	}
	out := make([]Info, 0, len(ids))
	for _, id := range ids {
		out = append(out, Info{
			ID:          id,
			Facing:      Facing(m.Facing(id)),
			Orientation: m.Orientation(id),
		})
	}
	return out, nil
}
