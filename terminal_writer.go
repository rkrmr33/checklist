package checklist

import (
	"fmt"
	"io"
)

type (
	TerminalWriter struct {
		io.Writer
	}
)

func NewTerminalWriter(w io.Writer) *TerminalWriter {
	return &TerminalWriter{
		Writer: w,
	}
}

func (w *TerminalWriter) Clean(lines int) error {
	if lines == 0 {
		return nil // do nothing
	}

	if lines < 0 {
		return w.clearScreen()
	}

	return w.clearLines(lines)
}

func (w *TerminalWriter) clearScreen() error {
	if _, err := fmt.Fprint(w, "\033[H\033[2J"); err != nil {
		return err
	}

	_, err := fmt.Fprint(w, "\033[0;0H")
	return err
}

func (w *TerminalWriter) clearLines(lines int) error {
	_, err := fmt.Fprintf(w, "\033[%dF\033[J", lines)
	return err
}
