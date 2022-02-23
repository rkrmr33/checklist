// +build windows

package checklist

import (
	"golang.org/x/sys/windows"
)

func (w *terminalWriter) getWidth() int {
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(w.fd), &info); err != nil {
		return defaultTerminalWidth
	}

	return int(info.Window.Right - info.Window.Left)
}
