package sdk

// Version is a dotted version out of the Android SDK — "37.0.0" for the build
// tools, "30.0.16138531" for the NDK — with the suffix a preview carries.
//
// The parts are kept as they were written rather than folded into three named
// fields, because the SDK is not consistent about how many there are.
type Version struct {
	Parts []int
	// Pre is what followed a hyphen: "rc2", "beta1". Empty on a release.
	Pre string
}
