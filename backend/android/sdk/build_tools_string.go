package sdk

// String is the build-tools version as a person would say it.
func (b BuildTools) String() string { return "build-tools " + b.Version.String() }
