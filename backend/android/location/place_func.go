//go:build android

package location

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func place(e *jni.Env, a jni.Object) (Place, error) {
	var p Place
	// An Address holds its lines separately and the whole thing nowhere, so
	// the one line everybody wants has to be put back together.
	last, err := e.InvokeInt(a, "getMaxAddressLineIndex", jni.Sig(jni.TInt))
	if err != nil {
		return p, err
	}
	for i := range int(last) + 1 {
		line, err := e.InvokeString(a, "getAddressLine",
			jni.Sig(jni.TString, jni.TInt), jni.Int(int32(i)))
		if err != nil || line == "" {
			continue
		}
		if p.Address != "" {
			p.Address += ", "
		}
		p.Address += line
	}
	for _, f := range []struct {
		method string
		into   *string
	}{
		{"getLocality", &p.Locality},
		{"getAdminArea", &p.Region},
		{"getCountryName", &p.Country},
		{"getPostalCode", &p.PostCode},
	} {
		v, err := e.InvokeString(a, f.method, jni.Sig(jni.TString))
		if err != nil {
			return p, err
		}
		*f.into = v
	}
	if p.Latitude, err = invokeDouble(e, a, "getLatitude"); err != nil {
		return p, err
	}
	if p.Longitude, err = invokeDouble(e, a, "getLongitude"); err != nil {
		return p, err
	}
	return p, nil
}
