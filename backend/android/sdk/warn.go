package sdk

func (d *Doctor) warn(text, fix string) {
	d.Findings = append(d.Findings, Finding{Warn, text, fix})
}
