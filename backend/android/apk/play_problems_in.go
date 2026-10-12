package apk

import (
	"archive/zip"
	"fmt"
	"strings"
)

// PlayProblemsIn is the same question asked of a package that has already
// been built — an .apk or an .aab — because a config is what was meant and
// the file is what will be uploaded.
//
// It reads the archive and does not run anything, so it works on a package
// built somewhere else.
func PlayProblemsIn(pkg string) ([]Problem, error) {
	z, err := zip.OpenReader(pkg)
	if err != nil {
		return nil, fmt.Errorf("apk: reading %s: %w", pkg, err)
	}
	defer z.Close()

	abis, compressed, symbols, size := scanPackage(&z.Reader)
	bundle := strings.HasSuffix(pkg, ".aab")
	return packageProblems(abis, compressed, symbols, bundle, size), nil
}
