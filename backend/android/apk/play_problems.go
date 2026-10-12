package apk

// PlayProblems is what the store will say about this config, answered before
// anything is built rather than after it is uploaded.
//
// It reads the config and nothing else, so it is worth calling early: the
// whole point is to fail in a second instead of after a two-minute build and
// a day of waiting.
func (c Config) PlayProblems() []Problem {
	// A copy, so that asking the question does not change the answer: check
	// fills in the defaults, and the defaults are most of what is being
	// judged.
	cfg := c
	if err := cfg.check(); err != nil {
		return []Problem{{Fatal: true, Text: err.Error()}}
	}
	var ps []Problem
	ps = playProblemsABIs(cfg, ps)
	ps = playProblemsStore(cfg, ps)
	ps = playProblemsOptional(cfg, ps)
	return ps
}
