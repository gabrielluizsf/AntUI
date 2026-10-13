package sdk

import "fmt"

// ParseABI reads an ABI name, also accepting the Go architecture that maps
// onto it so a caller can pass GOARCH through unchanged.
func ParseABI(s string) (ABI, error) {
	switch s {
	case "arm64-v8a", "arm64":
		return Arm64, nil
	case "armeabi-v7a", "armeabi", "arm":
		return Arm, nil
	case "x86_64", "amd64":
		return X64, nil
	case "x86", "386":
		return X86, nil
	}
	return "", fmt.Errorf("sdk: %q is not an Android ABI", s)
}
