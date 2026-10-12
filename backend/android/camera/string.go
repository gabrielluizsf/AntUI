//go:build android

package camera

// String names the direction the camera points.
func (f Facing) String() string {
	switch f {
	case Front:
		return "front"
	case Back:
		return "back"
	}
	return "external"
}
