package apk

// checkLinks refuses a link the platform would not accept.
func (c *Config) checkLinks() error {
	for _, l := range c.Links {
		if err := l.valid(); err != nil {
			return err
		}
	}
	return nil
}
