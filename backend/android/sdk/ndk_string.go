package sdk

// String is the NDK version as a person would say it.
func (n NDK) String() string { return "NDK " + n.Version.String() }
