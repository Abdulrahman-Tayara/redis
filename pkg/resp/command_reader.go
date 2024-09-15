package resp

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidRespCommand = errors.New("INVALID_RESP_COMMAND")
	ErrReaderRead         = errors.New("READ_ERROR")
)

type CommandReader interface {
	// ReadCommand returns the command, command's args and the error
	ReadCommand() (string, []any, error)
}

type commandReader struct {
	r io.Reader

	buf int
}

func NewCommandReader(r io.Reader, buf int) CommandReader {
	return &commandReader{
		r:   r,
		buf: buf,
	}
}

func (r *commandReader) ReadCommand() (string, []any, error) {
	buffer := make([]byte, r.buf)

	n, err := r.r.Read(buffer)
	if err != nil {
		return "", nil, fmt.Errorf("%w, err: %v", ErrReaderRead, err.Error())
	}

	respValue, err := Unmarshal(buffer[:n])
	if err != nil {
		return "", nil, err
	}

	switch respValue.(type) {
	case string:
		return respValue.(string), []any{}, nil
	case []any:
		arr := respValue.([]any)
		if len(arr) == 0 {
			return "", nil, ErrInvalidRespCommand
		}

		command, ok := arr[0].(string)
		if !ok {
			return "", nil, ErrInvalidRespCommand
		}

		return command, arr[1:], nil
	}

	return "", nil, ErrInvalidRespCommand
}
