package sdk

// BuildTools is one entry under the SDK's build-tools directory: aapt2, d8,
// zipalign and apksigner live here.
type BuildTools struct {
	Dir     string
	Version Version
}
