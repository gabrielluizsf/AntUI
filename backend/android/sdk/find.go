package sdk

import "strconv"

// Find locates everything Require asks for. The error it returns on a
// partial install is a [*Missing], which names the packages to install and
// the sdkmanager line that installs them.
func Find(req Require) (*Toolchain, error) {
	root := findRoot()
	if root == "" {
		return nil, ErrNoSDK
	}
	t := &Toolchain{Root: root}
	miss := &Missing{Root: root}

	if bt, ok := findBuildTools(root); !ok {
		miss.add("build-tools", "build-tools;37.0.0")
	} else if !bt.Version.AtLeast(req.MinBuildTools) {
		miss.old("build-tools", bt.Version.String(), req.MinBuildTools.String())
		t.BuildTools = bt
	} else {
		t.BuildTools = bt
	}

	if p, ok := findPlatform(root); !ok {
		miss.add("platforms", "platforms;android-37.2")
	} else if p.API < req.MinPlatform {
		miss.old("platforms", p.Name, "android-"+strconv.Itoa(req.MinPlatform))
		t.Platform = p
	} else {
		t.Platform = p
	}

	// Both of these are looked for whether or not they were asked for: a
	// caller that says it does not need a JDK still signs with one if there
	// is one, and finding out later that it was there all along is a worse
	// failure than the search costs.
	if n, ok := findNDK(root); ok {
		t.NDK = n
	} else if req.NDK {
		miss.add("ndk", "ndk;30.0.16138531")
	}
	if j, ok := findJDK(); ok {
		t.JDK = j
	} else if req.JDK {
		miss.add("a JDK", "")
	}
	t.Bundletool = findBundletool(root)
	if len(miss.Problems) > 0 {
		return t, miss
	}
	return t, nil
}
