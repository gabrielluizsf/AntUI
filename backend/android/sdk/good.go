package sdk

const (
	// Good is something found and usable.
	Good Level = iota
	// Warn is usable but not what to publish from.
	Warn
	// Bad stops a build.
	Bad
)
