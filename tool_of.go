//go:build android

package antui

import "github.com/gabrielluizsf/antui/backend/android/ndk"

func toolOf(t int) Tool {
	switch t {
	case ndk.ToolFinger:
		return ToolFinger
	case ndk.ToolStylus:
		return ToolStylus
	case ndk.ToolMouse:
		return ToolMouse
	case ndk.ToolEraser:
		return ToolEraser
	}
	return ToolUnknown
}
