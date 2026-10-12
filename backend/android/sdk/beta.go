package sdk

// Beta reports whether this NDK is a preview. Google publishes betas into the
// same channel as releases, so an install can quietly be one; a build still
// works, but it is not what to ship a store bundle from.
func (n NDK) Beta() bool { return !n.Version.Release() }
