//go:build android

package camera

// Default is the first camera pointing a given way.
func Default(f Facing) (Info, bool, error) {
	all, err := List()
	if err != nil {
		return Info{}, false, err
	}
	for _, c := range all {
		if c.Facing == f {
			return c, true, nil
		}
	}
	return Info{}, false, nil
}
