package css

import (
	"strings"
)

// supportsCondition reads the prelude of an @supports rule and answers it
// against the engine: whether it has the properties, the values and the
// selectors the condition tests. A condition the reader cannot make sense of
// is reported and answered false — the block it guards leaves the sheet, which
// is what a browser does with a test it cannot parse — while a test that
// simply fails is answered false in silence: that is the test working.
func supportsCondition(prelude string) (bool, []string) {
	r := &supportsReader{src: stripComments(prelude)}
	ok := r.condition()
	if r.bad || !r.atEnd() {
		return false, []string{fmtErrf("ignoring @supports condition %q", strings.TrimSpace(prelude)).Error()}
	}
	return ok, nil
}
