//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package antui

func (n *x11Native) sendNotification(title, body string) bool {
	return n.d.SendNotification(title, body)
}
