package checklist

import (
	"bytes"
	"io"
)

type (
	fileWriter struct {
		io.Writer
		buffer    *bytes.Buffer
		lastBytes []byte
	}
)

func newFileWriter(w io.Writer) *fileWriter {
	return &fileWriter{
		Writer:    w,
		buffer:    &bytes.Buffer{},
		lastBytes: make([]byte, 0),
	}
}

func (w *fileWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *fileWriter) clean(fullscreen bool) error {
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
