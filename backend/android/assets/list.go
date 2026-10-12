//go:build android

package assets

import (
	"io/fs"
	"sort"
)

func list(dir string) ([]fs.DirEntry, error) {
	names, err := listNames(dir)
	if err != nil || len(names) == 0 {
		return nil, err
	}
	out := make([]fs.DirEntry, 0, len(names))
	for _, name := range names {
		// A name whose own listing has anything in it is a directory. There
		// is no other way to ask: the platform reports a flat list of names
		// and nothing about what they are. An empty directory cannot be told
		// from a file this way, and does not need to be — aapt2 drops empty
		// directories on the way into the package.
		child := name
		if d := inner(dir); d != "" {
			child = d + "/" + name
		}
		sub, _ := listNames(child)
		out = append(out, entry{info: info{name: name, dir: len(sub) > 0}})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out, nil
}
