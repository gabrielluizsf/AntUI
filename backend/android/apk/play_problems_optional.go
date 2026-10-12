package apk

// playProblemsOptional judges the things that cost without failing: an icon,
// the symbols left in the library, and the version code.
func playProblemsOptional(cfg Config, ps []Problem) []Problem {
	if cfg.Icon == "" {
		ps = append(ps, Problem{false,
			"no icon",
			"the app appears under the platform's blank one, and the listing " +
				"needs a 512x512 one separately"})
	}
	if cfg.Strip != nil && !*cfg.Strip {
		ps = append(ps, Problem{false,
			"the library ships with its symbols in it",
			"it is about a third larger than it needs to be, and the store " +
				"asks for the symbols as a separate file anyway"})
	}
	if cfg.VersionCode <= 0 {
		ps = append(ps, Problem{true,
			"no version code",
			"it must be a positive number and must rise with every upload; a " +
				"code that has been used once can never be used again"})
	}
	return ps
}
