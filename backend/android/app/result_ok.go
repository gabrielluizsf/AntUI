//go:build android

package app

// The two result codes android.app.Activity defines. RESULT_OK is -1 and not
// 0, which is the opposite of every other convention and is worth stating.
const (
	ResultOK       = -1
	ResultCanceled = 0
)
