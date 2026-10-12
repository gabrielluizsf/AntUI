package sdk

// Finding is one line of the report.
type Finding struct {
	Level Level
	Text  string
	// Fix is what to run, when running something would help.
	Fix string
}
