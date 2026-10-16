package svg

// fetchesOffFrame is whether an address is read away from the frame that asked
// for it. A server is: waiting on one inside a frame is the wait this keeps out
// of the drawing, so its picture arrives in a later frame instead of holding
// the frame it was named in. A file on this machine and the bytes a `data:`
// address already carries are not — they answer in the time the read takes,
// and a drawing that named one has it to draw in the same frame.
func fetchesOffFrame(target string) bool {
	switch hrefScheme(target) {
	case "http", "https":
		return true
	}
	return false
}
