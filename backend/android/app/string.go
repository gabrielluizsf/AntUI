//go:build android

package app

// String names the event, which is what a log line wants.
func (e Event) String() string {
	if int(e) < len(eventNames) {
		return eventNames[e]
	}
	return "Event(?)"
}
