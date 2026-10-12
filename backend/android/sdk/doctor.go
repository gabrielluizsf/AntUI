package sdk

// Doctor is what is on this machine and what is wrong with it. It never
// fails: a missing SDK is a finding, not an error, because the whole point
// is to be able to say so.
type Doctor struct {
	Toolchain *Toolchain // nil when there is no SDK at all
	Findings  []Finding
}
