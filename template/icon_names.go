package template

import "sort"

// IconNames are the icons registered right now, sorted, which is what a tool
// that lets a developer pick an icon wants to show.
func IconNames() []string {
	iconMu.RLock()
	defer iconMu.RUnlock()
	names := make([]string, 0, len(icons))
	for name := range icons {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
