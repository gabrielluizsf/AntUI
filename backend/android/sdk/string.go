package sdk

// String names how bad a finding is.
func (l Level) String() string {
	switch l {
	case Good:
		return "ok"
	case Warn:
		return "warn"
	}
	return "missing"
}
