package sdk

func (d *Doctor) good(what, detail string) {
	d.Findings = append(d.Findings, Finding{Good, what + ": " + detail, ""})
}
