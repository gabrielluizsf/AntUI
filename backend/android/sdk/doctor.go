package sdk

import (
	"fmt"
	"runtime"
	"strings"
)

// Doctor is what is on this machine and what is wrong with it. It never
// fails: a missing SDK is a finding, not an error, because the whole point
// is to be able to say so.
type Doctor struct {
	Toolchain *Toolchain // nil when there is no SDK at all
	Findings  []Finding
}

// Finding is one line of the report.
type Finding struct {
	Level Level
	Text  string
	// Fix is what to run, when running something would help.
	Fix string
}

// Level is how much a finding matters.
type Level int

const (
	// Good is something found and usable.
	Good Level = iota
	// Warn is usable but not what to publish from.
	Warn
	// Bad stops a build.
	Bad
)

// String names how bad a finding is.
func (l Level) String() string {
	switch l {
	case Good:
		return "ok"
	case Warn:
		return "warn"
	}
	return "missing"
}

// Check looks the machine over. Req says what to complain about the absence
// of; everything found is reported either way.
func Check(req Require) *Doctor {
	d := &Doctor{}
	t, err := Find(req)
	d.Toolchain = t
	if t == nil {
		d.bad("no Android SDK found",
			"install the command line tools, or set ANDROID_HOME to an SDK you already have")
		return d
	}
	d.good("SDK", t.Root)

	if t.BuildTools.Dir != "" {
		d.good("build-tools", t.BuildTools.Version.String())
	}
	if t.Platform.Dir != "" {
		d.good("platform", t.Platform.String())
	}
	if t.NDK.Dir != "" {
		d.good("NDK", t.NDK.Version.String()+" ("+t.NDK.Host+")")
		if t.NDK.Beta() {
			d.warn("the NDK is a preview: "+t.NDK.Version.String(),
				"a build works, but install a release NDK before making a bundle to upload")
		}
		if _, err := t.NDK.Clang(Arm64, MinSDK); err != nil {
			d.bad(err.Error(), "install an NDK that still supports API "+itoa(MinSDK))
		}
	}
	if t.JDK.Javac != "" {
		d.good("JDK", t.JDK.Dir)
	}
	if exists(t.Adb()) {
		d.good("adb", t.Adb())
	} else {
		d.warn("no adb", `sdkmanager "platform-tools"`)
	}
	if t.Bundletool != "" {
		d.good("bundletool", t.Bundletool)
	} else {
		// Not part of the SDK and not needed to build anything, so this is a
		// note and not a warning: only checking a bundle wants it.
		d.good("bundletool", "not here — a bundle can be built without it, "+
			"but not taken apart again")
	}
	if t.Platform.API < PlayTarget {
		d.warn(fmt.Sprintf("the newest platform is %s, and the Play Store wants an app to target API %d or newer",
			t.Platform, PlayTarget),
			fmt.Sprintf(`sdkmanager "platforms;android-%d"`, PlayTarget))
	}

	// Everything Find complained about, in its own words.
	var miss *Missing
	if as(err, &miss) {
		for _, p := range miss.Problems {
			if p.Have != "" {
				d.bad(p.What+" is "+p.Have+", want "+p.Want+" or newer", "")
			} else {
				d.bad(p.What+" is not installed", fixLine(p.Install))
			}
		}
	}
	return d
}

// MinSDK is the oldest Android this library supports: API 21, Android 5.0.
// It is not an arbitrary floor — it is the oldest the current NDK will build
// for, and below it there is no compiler to use.
const MinSDK = 21

// PlayTarget is the API level the Play Store requires a new upload to target.
// It moves every year, one year behind the newest Android.
const PlayTarget = 36

func (d *Doctor) good(what, detail string) {
	d.Findings = append(d.Findings, Finding{Good, what + ": " + detail, ""})
}
func (d *Doctor) warn(text, fix string) {
	d.Findings = append(d.Findings, Finding{Warn, text, fix})
}
func (d *Doctor) bad(text, fix string) {
	d.Findings = append(d.Findings, Finding{Bad, text, fix})
}

// OK reports whether anything found would stop a build.
func (d *Doctor) OK() bool {
	for _, f := range d.Findings {
		if f.Level == Bad {
			return false
		}
	}
	return true
}

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

func fixLine(pkg string) string {
	if pkg == "" {
		return ""
	}
	return `sdkmanager "` + pkg + `"`
}

func itoa(n int) string { return fmt.Sprintf("%d", n) }
