package font

import "fmt"

// readCmapSubtable prepares a cmap subtable: the format, the data, and the
// segment or group count its own format carries.
func readCmapSubtable(sub []byte) (CmapTable, error) {
	t := CmapTable{Format: be16(sub, 0), Data: sub}
	switch t.Format {
	case 4:
		if len(sub) < 14 {
			return CmapTable{}, fmt.Errorf("antUI: a truncated format 4 cmap")
		}
		t.Segments = int(be16(sub, 6)) / 2
	case 12:
		if len(sub) < 16 {
			return CmapTable{}, fmt.Errorf("antUI: a truncated format 12 cmap")
		}
		t.Groups = int(be32(sub, 12))
	default:
		return CmapTable{}, fmt.Errorf("antUI: cmap format %d is not read here", t.Format)
	}
	return t, nil
}
