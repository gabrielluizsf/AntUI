package sdk

// GOARCH is the Go architecture that builds for this ABI.
func (a ABI) GOARCH() string {
	switch a {
	case Arm64:
		return "arm64"
	case Arm:
		return "arm"
	case X64:
		return "amd64"
	case X86:
		return "386"
	}
	return ""
}
