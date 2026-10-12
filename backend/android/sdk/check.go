package sdk

// Check looks the machine over. Req says what to complain about the absence
// of; everything found is reported either way.
func Check(req Require) *Doctor {
	d := &Doctor{}
	t, err := Find(req)
	d.Toolchain = t
	if t == nil {
		d.bad("no Android SDK found",
			"install the command line tools, or set ANDROID_HOME to an SDK you already have")
		return d
	}
	d.good("SDK", t.Root)
	reportFound(d, t)
	reportMissing(d, err)
	return d
}
