package checklist

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	defaultTerminalWidth = int(math.MaxInt32)
)

type (
	terminalWriter struct {
		io.Writer
		fd        int
		buffer    *bytes.Buffer
		lastBytes []byte
	}
)

func newTerminalWriter(w io.Writer, fd int) *terminalWriter {
	return &terminalWriter{
		Writer:    w,
		fd:        fd,
		buffer:    &bytes.Buffer{},
		lastBytes: make([]byte, 0),
	}
}

func (w *terminalWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *terminalWriter) flush(fullscreen bool) error {
	defer w.buffer.Reset()

	if w.buffer.String() == string(w.lastBytes) {
		return nil
	}

	newLastBytes := make([]byte, w.buffer.Len())
	copy(newLastBytes, w.buffer.Bytes())

	// clean last printed message
	if err := w.clean(w.calcLines(w.lastBytes, fullscreen)); err != nil {
		return err
	}

	if _, err := io.Copy(w.Writer, w.buffer); err != nil {
		return err
	}

	w.lastBytes = newLastBytes

	return nil
}

func (w *terminalWriter) clean(lines int) error {
	if lines == 0 {
		return nil // do nothing
	}

	if lines < 0 {
		return w.clearScreen()
	}

	return w.clearLines(lines)
}

func (w *terminalWriter) clearScreen() error {
	if err := w.deleteAllLines(); err != nil {
		return err
	}

	return w.moveTopLeft()
}

func (w *terminalWriter) clearLines(lines int) error {
	if err := w.moveUpLines(lines); err != nil {
		return err
	}

	return w.clearFromCursorToEnd()
}

func (w *terminalWriter) calcLines(data []byte, fullscreen bool) int {
	if fullscreen {
		return -1 // will clear all screen
	}

	screenWidth := w.getWidth()
	rawLines := strings.Split(string(data), "\n")
	n := len(rawLines) - 1

	for _, line := range rawLines {
		n += runewidth.StringWidth(stripANSI(line)) / screenWidth
	}

	return n
}

func (w *terminalWriter) moveUpLines(lines int) error {
	_, err := fmt.Fprintf(w.Writer, "\033[%dA\r", lines)
	return err
}

func (w *terminalWriter) clearFromCursorToEnd() error {
	_, err := fmt.Fprint(w.Writer, "\033[J")
	return err
}

func (w *terminalWriter) deleteAllLines() error {
	_, err := fmt.Fprint(w.Writer, "\033[H\033[2J")
	return err
}

func (w *terminalWriter) moveTopLeft() error {
	_, err := fmt.Fprint(w.Writer, "\033[0;0H")
	return err
}
