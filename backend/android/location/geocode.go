//go:build android

package location

import (
	"errors"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Place is an address, and where it is.
type Place struct {
	// Address is the whole thing as one line, the way it would be written on
	// an envelope in that country.
	Address string
	// Locality is the town or city, Region the state or province.
	Locality string
	Region   string
	Country  string
	PostCode string
	// Latitude and Longitude are where the address is, which for a forward
	// lookup is the answer and for a reverse one is roughly what was asked.
	Latitude, Longitude float64
}

// ErrNoGeocoder is a device with no service behind the geocoder. It is not
// rare: the geocoder is a network service and a device without Play Services
// or without a connection has none.
var ErrNoGeocoder = errors.New("location: this device has no geocoder")

// Available reports whether geocoding will work at all.
func Available() bool {
	var out bool
	jni.Do(func(e *jni.Env) error {
		var err error
		c, err := e.Class("android/location/Geocoder")
		if err != nil {
			return err
		}
		m, err := e.StaticMethod(c, "isPresent", jni.Sig(jni.TBool))
		if err != nil {
			return err
		}
		out, err = e.CallStaticBool(c, m)
		return err
	})
	return out
}

// Geocode turns an address into places. max bounds how many come back; one
// is usually right, and more is for offering a choice.
//
// **It goes to the network and blocks.** Not from the frame loop, and not
// from anything holding a lock a frame wants.
func Geocode(query string, max int) ([]Place, error) {
	if max <= 0 {
		max = 1
	}
	var out []Place
	err := jni.Do(func(e *jni.Env) error {
		coder, err := geocoder(e)
		if err != nil {
			return err
		}
		q, err := e.String(query)
		if err != nil {
			return err
		}
		list, err := e.Invoke(coder, "getFromLocationName",
			jni.Sig(jni.TClass("java/util/List"), jni.TString, jni.TInt),
			jni.Ref(q), jni.Int(int32(max)))
		if err != nil {
			return err
		}
		out, err = places(e, list)
		return err
	})
	return out, err
}

// Reverse turns a position into addresses.
//
// The same warning: it goes to the network and blocks.
func Reverse(latitude, longitude float64, max int) ([]Place, error) {
	if max <= 0 {
		max = 1
	}
	var out []Place
	err := jni.Do(func(e *jni.Env) error {
		coder, err := geocoder(e)
		if err != nil {
			return err
		}
		list, err := e.Invoke(coder, "getFromLocation",
			jni.Sig(jni.TClass("java/util/List"), jni.TDouble, jni.TDouble, jni.TInt),
			jni.Double(latitude), jni.Double(longitude), jni.Int(int32(max)))
		if err != nil {
			return err
		}
		out, err = places(e, list)
		return err
	})
	return out, err
}

func geocoder(e *jni.Env) (jni.Object, error) {
	return e.Make("android/location/Geocoder",
		jni.Sig(jni.TVoid, jni.TContext), jni.Ref(app.Context()))
}

// places reads a java.util.List of Address into Go.
func places(e *jni.Env, list jni.Object) ([]Place, error) {
	if list.IsNil() {
		return nil, nil
	}
	n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
	if err != nil || n == 0 {
		return nil, err
	}
	out := make([]Place, 0, n)
	err = e.Frame(int(n)*8+16, func() error {
		for i := range int(n) {
			item, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt), jni.Int(int32(i)))
			if err != nil || item.IsNil() {
				continue
			}
			p, err := place(e, item)
			if err != nil {
				return err
			}
			out = append(out, p)
		}
		return nil
	})
	return out, err
}

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
