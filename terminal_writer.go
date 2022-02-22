package checklist

import (
	"fmt"
	"io"
	"sync"
)

type (
	TerminalWriter struct {
		io.Writer
		initialCursorPos sync.Once
	}
)

func NewTerminalWriter(w io.Writer) *TerminalWriter {
	return &TerminalWriter{
		Writer: w,
	}
}

func (w *TerminalWriter) Write(data []byte) (int, error) {
	var err error

	w.initialCursorPos.Do(func() {
		err = w.saveCursorPos()
	})
	if err != nil {
		return 0, err
	}

	return w.Write(data)
}

func (w *TerminalWriter) Clean(fullscreen bool) error {
	if fullscreen {
		return w.clearScreen()
	}

	return w.clearLines()
}

func (w *TerminalWriter) clearScreen() error {
	if err := w.deleteAllLines(); err != nil {
		return err
	}

	return w.moveTopLeft()
}

func (w *TerminalWriter) clearLines() error {
	if err := w.restoreCursorPos(); err != nil {
		return err
	}

	if err := w.clearFromCursorToEnd(); err != nil {
		return err
	}

	return w.saveCursorPos()
}

func (w *TerminalWriter) saveCursorPos() error {
	_, err := fmt.Fprint(w, "\033[s")
	return err
}

func (w *TerminalWriter) restoreCursorPos() error {
	_, err := fmt.Fprint(w, "\033[u")
	return err
}

func (w *TerminalWriter) clearFromCursorToEnd() error {
	_, err := fmt.Fprint(w, "\033[J")
	return err
}

func (w *TerminalWriter) deleteAllLines() error {
	_, err := fmt.Fprint(w, "\033[H\033[2J")
	return err
}

func (w *TerminalWriter) moveTopLeft() error {
	_, err := fmt.Fprint(w, "\033[0;0H")
	return err
}
