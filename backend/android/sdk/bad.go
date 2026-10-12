package sdk

func (d *Doctor) bad(text, fix string) {
	d.Findings = append(d.Findings, Finding{Bad, text, fix})
}
