package checklist

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type (
	terminalWriter struct {
		io.Writer
		initialCursorPos sync.Once
		buffer           *bytes.Buffer
		lastBytes        []byte
	}
)

func newTerminalWriter(w io.Writer) *terminalWriter {
	return &terminalWriter{
		Writer:           w,
		initialCursorPos: sync.Once{},
		buffer:           &bytes.Buffer{},
		lastBytes:        make([]byte, 0),
	}
}

func (w *terminalWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *terminalWriter) flush() error {
	defer w.buffer.Reset()

	if w.buffer.String() == string(w.lastBytes) {
		return nil
	}

	newLastBytes := make([]byte, w.buffer.Len())
	copy(newLastBytes, w.buffer.Bytes())
	w.lastBytes = newLastBytes

	if _, err := io.Copy(w.Writer, w.buffer); err != nil {
		return err
	}

	return nil
}

func (w *terminalWriter) clean(fullscreen bool) error {
	var err error

	w.initialCursorPos.Do(func() {
		fmt.Printf("h")
		err = w.saveCursorPos()
	})
	if err != nil {
		return err
	}

	if fullscreen {
		w.clearScreen()
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
	_, err := fmt.Fprint(w.buffer, "\033[s")
	return err
}

func (w *terminalWriter) restoreCursorPos() error {
	_, err := fmt.Fprint(w.buffer, "\033[u")
	return err
}

func (w *terminalWriter) clearFromCursorToEnd() error {
	_, err := fmt.Fprint(w.buffer, "\033[J")
	return err
}

func (w *terminalWriter) deleteAllLines() error {
	_, err := fmt.Fprint(w.buffer, "\033[H\033[2J")
	return err
}

func (w *terminalWriter) moveTopLeft() error {
	_, err := fmt.Fprint(w.buffer, "\033[0;0H")
	return err
}
