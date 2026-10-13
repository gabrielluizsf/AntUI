//go:build android

package picker

import "github.com/gabrielluizsf/antui/backend/android/app"

// Many asks for several at once. It comes back with one address per file, in
// the order the picker gave them.
func Many(mime string) ([]string, error) {
	in, err := build(mime, true)
	if err != nil {
		return nil, err
	}
	r, err := app.StartForResult(in)
	if err != nil {
		return nil, err
	}
	defer drop(r.Data)
	if !r.OK() {
		return nil, ErrCancelled
	}
	return addresses(r.Data)
}
