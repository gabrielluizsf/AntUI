//go:build android

package dialog

// The interfaces the dialog answers through. All three are interfaces, so
// one proxy implements the lot and no Java had to be written for any of them.
var listenerInterfaces = []string{
	"android.content.DialogInterface$OnClickListener",
	"android.content.DialogInterface$OnCancelListener",
	"android.content.DialogInterface$OnDismissListener",
}
