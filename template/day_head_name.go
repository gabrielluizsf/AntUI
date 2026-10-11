package template

func dayHeadName(wd int) string {
	return "SMTWTFS"[wd : wd+1]
}
