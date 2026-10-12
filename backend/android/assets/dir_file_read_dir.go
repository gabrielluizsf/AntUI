//go:build android

package assets

import (
	"io"
	"io/fs"
)

func (d *dirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if n <= 0 {
		out := d.entries[d.at:]
		d.at = len(d.entries)
		return out, nil
	}
	if d.at >= len(d.entries) {
		return nil, io.EOF
	}
	end := min(d.at+n, len(d.entries))
	out := d.entries[d.at:end]
	d.at = end
	return out, nil
}
