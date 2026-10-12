package sdk

import (
	"fmt"
	"runtime"
	"strings"
)

// String is the report, as it is meant to be printed.
func (d *Doctor) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "AntUI android, on %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
	for _, f := range d.Findings {
		fmt.Fprintf(&b, "  %-8s %s\n", f.Level, f.Text)
		if f.Fix != "" {
			fmt.Fprintf(&b, "           %s\n", f.Fix)
		}
	}
	b.WriteByte('\n')
	if d.OK() {
		b.WriteString("  ready to build.\n")
	} else {
		b.WriteString("  not ready.\n")
	}
	return b.String()
}
