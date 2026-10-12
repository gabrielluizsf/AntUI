package sdk

// String is the ABI as Android spells it, which is also the directory
// name inside a package.
func (a ABI) String() string { return string(a) }
