package font

import "fmt"

func parseTTFAt(data []byte, base int) (*TTF, error) {
	if base+12 > len(data) {
		return nil, fmt.Errorf("antUI: a truncated font")
	}
	tables, err := readTableDir(data, base)
	if err != nil {
		return nil, err
	}
	f := &TTF{Data: data}
	if err := readTTFTables(f, tables); err != nil {
		return nil, err
	}
	return f, nil
}
