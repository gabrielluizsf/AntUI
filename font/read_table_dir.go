package font

import "fmt"

// readTableDir maps each table tag in a font's directory to its bytes.
//
// A table that runs off the end is skipped rather than fatal: some fonts pad
// the last one, and refusing the whole file for a table nothing here reads
// would be absurd.
func readTableDir(data []byte, base int) (map[string][]byte, error) {
	count := int(be16(data, base+4))
	tables := map[string][]byte{}
	for i := range count {
		rec := base + 12 + i*16
		if rec+16 > len(data) {
			return nil, fmt.Errorf("antUI: a truncated table directory")
		}
		name := string(data[rec : rec+4])
		off, length := int(be32(data, rec+8)), int(be32(data, rec+12))
		if off < 0 || length < 0 || off+length > len(data) {
			continue
		}
		tables[name] = data[off : off+length]
	}
	return tables, nil
}
