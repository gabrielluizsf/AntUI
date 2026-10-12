package apk

// DeepLink is an address this app opens.
//
//	DeepLink{Scheme: "myapp"}                  myapp://anything
//	DeepLink{Scheme: "https", Host: "ex.com"}  every page on ex.com
//	DeepLink{Scheme: "https", Host: "ex.com", Path: "/open"}
type DeepLink struct {
	// Scheme is "https", or a scheme of the app's own. A custom scheme is
	// simpler and is claimed by whoever asks — any other app may claim the
	// same one, and the user is asked which. A web address cannot be stolen
	// that way, but has to be proved.
	Scheme string `json:"scheme"`
	// Host is the domain, for a web address.
	Host string `json:"host,omitempty"`
	// Path is a prefix, and empty means the whole host.
	Path string `json:"path,omitempty"`
	// Verify makes it an App Link: Android checks
	// https://<host>/.well-known/assetlinks.json for this app's signature
	// and, if it finds it, opens the address here **without asking the
	// user**. Without that file the filter still works and the user is asked.
	Verify bool `json:"verify,omitempty"`
}
