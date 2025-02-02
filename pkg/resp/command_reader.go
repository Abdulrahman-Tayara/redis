package resp

import (
	"bufio"
	"errors"
	"io"
)

var (
	ErrInvalidRespCommand = errors.New("INVALID_RESP_COMMAND")
)

type CommandReader interface {
	// ReadCommand returns the command, command's args and the error
	ReadCommand() (string, []any, error)
}

type commandReader struct {
	r            io.Reader
	bufferReader *bufio.Reader

	bufferedCommands []any
}

func NewCommandReader(r io.Reader) CommandReader {
	return &commandReader{
		r:            r,
		bufferReader: bufio.NewReader(r),
	}
}

func (r *commandReader) ReadCommand() (string, []any, error) {
	return r.readCommand()
}

func (r *commandReader) readCommand() (string, []any, error) {
	if len(r.bufferedCommands) > 0 {
		respSegment := r.bufferedCommands[0]
		r.bufferedCommands = r.bufferedCommands[1:]
		return parseRespCommand(respSegment)
	}

	respSegments, err := Unmarshal(r.bufferReader)
	if err != nil {
		return "", nil, err
	}
	if len(respSegments) == 0 {
		return "", nil, nil
	}
	respSegment := respSegments[0]
	r.bufferedCommands = respSegments[1:]

	return parseRespCommand(respSegment)
}

func parseRespCommand(respCommand any) (string, []any, error) {
	switch segment := respCommand.(type) {
	case string:
		return segment, []any{}, nil
	case []any:
		arr := segment
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
