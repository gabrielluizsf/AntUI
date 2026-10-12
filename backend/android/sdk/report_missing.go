package sdk

// reportMissing repeats what Find complained about, in its own words.
func reportMissing(d *Doctor, err error) {
	var miss *Missing
	if !as(err, &miss) {
		return
	}
	for _, p := range miss.Problems {
		if p.Have != "" {
			d.bad(p.What+" is "+p.Have+", want "+p.Want+" or newer", "")
		} else {
			d.bad(p.What+" is not installed", fixLine(p.Install))
		}
	}
}
