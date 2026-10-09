package css

import (
	"fmt"
)

func fmtErrf(format string, args ...any) error {
	return fmt.Errorf("css: "+format, args...)
}
