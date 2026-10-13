package sdk

import "fmt"

// reportFound says what the toolchain has, and what that means for a build.
func reportFound(d *Doctor, t *Toolchain) {
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
}
