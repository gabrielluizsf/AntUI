//go:build android

package app

// String names the state, which is what a log line wants.
func (s State) String() string {
	if int(s) < len(stateNames) {
		return stateNames[s]
	}
	return "state(?)"
}
