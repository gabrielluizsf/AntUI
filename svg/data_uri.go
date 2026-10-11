package svg

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// dataURI is the bytes a `data:` address carries: the media type it announces
// is left for the decoder to have an opinion about, base64 is unwrapped when
// the address says so — with or without the padding some encoders leave off —
// and anything else is the text it was written as, with what `%20` stands for
// put back.
func dataURI(href string) ([]byte, error) {
	meta, body, ok := strings.Cut(strings.TrimPrefix(href, "data:"), ",")
	if !ok {
		return nil, fmt.Errorf("the data: address has no comma in it")
	}
	for _, part := range strings.Split(meta, ";") {
		if strings.EqualFold(strings.TrimSpace(part), "base64") {
			clean := strings.Map(func(r rune) rune {
				if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
					return -1
				}
				return r
			}, body)
			if raw, err := base64.StdEncoding.DecodeString(clean); err == nil {
				return raw, nil
			}
			return base64.RawStdEncoding.DecodeString(clean)
		}
	}
	s, err := url.PathUnescape(body)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}
