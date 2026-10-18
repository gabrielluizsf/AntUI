package svg

import "strings"

// customProps picks the custom-property declarations out of an element's
// settled attributes: the `--name` ones are values for var() to find rather
// than attributes of the element, and they inherit instead of being read off
// where they are written.
func customProps(attrs map[string]string) map[string]string {
	var out map[string]string
	for name, raw := range attrs {
		if !strings.HasPrefix(name, "--") {
			continue
		}
		if out == nil {
			out = map[string]string{}
		}
		out[name] = raw
	}
	return out
}
