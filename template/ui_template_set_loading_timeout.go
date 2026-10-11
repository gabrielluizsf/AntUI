package template

import "time"

func (t *uiTemplate) SetLoadingTimeout(d time.Duration) { t.ld.setTimeout(d) }
