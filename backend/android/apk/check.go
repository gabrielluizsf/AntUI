package apk

// check fills in the defaults and reports what is still missing. It is
// called by every entry point, so a caller who skips Build still gets it.
func (c *Config) check() error {
	if err := c.checkRequired(); err != nil {
		return err
	}
	if err := c.checkVersions(); err != nil {
		return err
	}
	if err := c.checkABIs(); err != nil {
		return err
	}
	if err := c.checkLinks(); err != nil {
		return err
	}
	if err := c.checkOptional(); err != nil {
		return err
	}
	c.Keystore.fill()
	return nil
}
