package sdk

func fixLine(pkg string) string {
	if pkg == "" {
		return ""
	}
	return `sdkmanager "` + pkg + `"`
}
