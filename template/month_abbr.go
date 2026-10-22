package template

// monthAbbr is the month's name cut to three letters, which is what fits a
// cell in the month grid that [monthName]'s full word does not.
func monthAbbr(m int) string {
	name := monthName(m)
	if len(name) <= 3 {
		return name
	}
	return name[:3]
}
