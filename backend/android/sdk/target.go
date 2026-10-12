package sdk

// Target is the API level to write into a manifest as targetSdkVersion. It
// is the major number: the minor of an android-36.1 is a mid-year addition
// to the same level and is not what targetSdkVersion holds.
func (p Platform) Target() int { return p.API }
