package checklist

import (
	"fmt"
	"io"
	"sync"
)

type (
	terminalWriter struct {
		io.Writer
		initialCursorPos sync.Once
	}
)

func newTerminalWriter(w io.Writer) *terminalWriter {
	return &terminalWriter{
		Writer: w,
	}
}

func (w *terminalWriter) Write(data []byte) (int, error) {
	var err error

	w.initialCursorPos.Do(func() {
		err = w.saveCursorPos()
	})
	if err != nil {
		return 0, err
	}

	return w.Writer.Write(data)
}

func (w *terminalWriter) clean(fullscreen bool) error {
	if fullscreen {
		return w.clearScreen()
	}

	return w.clearLines()
}

func (w *terminalWriter) clearScreen() error {
	if err := w.deleteAllLines(); err != nil {
		return err
	}

	return w.moveTopLeft()
}

func (w *terminalWriter) clearLines() error {
	if err := w.restoreCursorPos(); err != nil {
		return err
	}

	if err := w.clearFromCursorToEnd(); err != nil {
		return err
	}

	return w.saveCursorPos()
}

func (w *terminalWriter) saveCursorPos() error {
	_, err := fmt.Fprint(w.Writer, "\033[s")
	return err
}

func (w *terminalWriter) restoreCursorPos() error {
	_, err := fmt.Fprint(w, "\033[u")
	return err
}

func (w *terminalWriter) clearFromCursorToEnd() error {
	_, err := fmt.Fprint(w, "\033[J")
	return err
}

func (w *terminalWriter) deleteAllLines() error {
	_, err := fmt.Fprint(w, "\033[H\033[2J")
	return err
}

func (w *terminalWriter) moveTopLeft() error {
	_, err := fmt.Fprint(w, "\033[0;0H")
	return err
}
