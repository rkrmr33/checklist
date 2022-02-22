// +build !windows
package checklist

import (
	"golang.org/x/sys/unix"
)

func (w *terminalWriter) getWidth() int {
	ws, err := unix.IoctlGetWinsize(w.fd, unix.TIOCGWINSZ)
	if err != nil {
		return defaultTerminalWidth
	}

	return int(ws.Col)
}
