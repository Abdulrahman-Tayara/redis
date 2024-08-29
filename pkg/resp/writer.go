package resp

import (
	"io"
)

type Writer struct {
	w io.Writer
}

func NewRespWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) WriteAny(v any) (n int, err error) {
	s, err := Marshal(v)
	if err != nil {
		return 0, err
	}

	return w.w.Write([]byte(s))
}
