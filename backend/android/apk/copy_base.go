package apk

import "archive/zip"

// copyBase copies every entry of the base APK, keeping whatever compression
// aapt2 chose — except resources.arsc, which since API 30 must also be
// stored, and is decompressed here if it is not.
func copyBase(zw *zip.Writer, zr *zip.Reader) error {
	for _, e := range zr.File {
		if e.Name == "resources.arsc" && e.Method != zip.Store {
			if err := storeEntry(zw, e); err != nil {
				return err
			}
			continue
		}
		if err := copyRaw(zw, e); err != nil {
			return err
		}
	}
	return nil
}
