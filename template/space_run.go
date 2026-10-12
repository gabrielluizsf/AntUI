package template

// spaceRun measures a run of soft spaces starting at i and returns its width
// and the index just past it. When spaces collapse the run counts as one.
func spaceRun(text string, i int, keepSpaces bool, spaceW int) (runW, next int) {
	j := i
	for j < len(text) && text[j] != '\n' && isSoftSpace(text[j]) {
		if keepSpaces {
			runW += spaceW
		}
		j++
	}
	if !keepSpaces {
		runW = spaceW
	}
	return runW, j
}
