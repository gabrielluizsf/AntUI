package sdk

// Avdmanager is what makes and lists the virtual devices.
func (t *Toolchain) Avdmanager() string { return t.cmdlineTool("avdmanager") }
