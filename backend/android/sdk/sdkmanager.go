package sdk

// Sdkmanager is what installs and updates everything else.
func (t *Toolchain) Sdkmanager() string { return t.cmdlineTool("sdkmanager") }
