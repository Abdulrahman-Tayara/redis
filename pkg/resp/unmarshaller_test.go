package resp

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshaller_UnmarshalString(t *testing.T) {
	tests := []unmarshalTest[string]{
		{
			"Valid simple string",
			stringReader("+OK\r\n"),
			"OK",
			nil,
		},
		{
			"Invalid simple string - missing CRLF",
			stringReader("+OK"),
			"",
			ErrValueNotEndWithCrlf,
		},
		{
			"Invalid simple string - contains LF",
			stringReader("+OK\nOK2\r\n"),
			"",
			errors.New("simple string mustn't contain a CR or LF character"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

func TestUnmarshaller_UnmarshalError(t *testing.T) {
	tests := []unmarshalTest[string]{
		{
			"Valid error",
			stringReader("-Error message\r\n"),
			"Error message",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

func TestUnmarshaller_UnmarshalInt(t *testing.T) {
	tests := []unmarshalTest[int64]{
		{
			"Valid int",
			stringReader(":10\r\n"),
			int64(10),
			nil,
		},
		{
			"Valid signed int",
			stringReader(":-12\r\n"),
			int64(-12),
			nil,
		},
		{
			"Invalid int",
			stringReader(":-12E\r\n"),
			0,
			errors.New("invalid int parsing"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

func TestUnmarshaller_UnmarshalBulkStrings(t *testing.T) {
	tests := []unmarshalTest[string]{
		{
			"Valid bulk strings",
			stringReader("$16\r\nHello\rMy\nFriends\r\n"),
			"Hello\rMy\nFriends",
			nil,
		},
		{
			"Valid empty bulk strings",
			stringReader("$0\r\n\r\n"),
			"",
			nil,
		},
		{
			"Invalid bulk strings - invalid string length",
			stringReader("$Hello\r\nWorld\r\n"),
			"",
			errors.New("invalid string length value, expected one or more decimal digits (0..9)"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

func TestUnmarshaller_UnmarshalArray(t *testing.T) {
	tests := []unmarshalTest[[]any]{
		{
			"Valid array",
			stringReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			[]any{"hello", "world"},
			nil,
		},
		{
			"Valid empty array",
			stringReader("*0\r\n"),
			[]any{},
			nil,
		},
		{
			"Valid multi-type array",
			stringReader("*4\r\n$5\r\nhello\r\n:12\r\n-Error Message\r\n+OK\r\n"),
			[]any{"hello", int64(12), "Error Message", "OK"},
			nil,
		},
		{
			"Valid nested arrays",
			stringReader("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n"),
			[]any{[]any{int64(1), int64(2), int64(3)}, []any{"Hello", "World"}},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

type unmarshalTest[T any] struct {
	name        string
	input       io.Reader
	expected    T
	expectedErr error
}

func unmarshalTestRunner[T any](test *unmarshalTest[T]) func(t *testing.T) {
	return func(t *testing.T) {
		output, err := Unmarshal(bufio.NewReader(test.input))

		assert.Equal(t, test.expectedErr, err)

		if err == nil {
			assert.Equal(t, test.expected, output[0])
		}
	}
}

func stringReader(s string) io.Reader {
	return strings.NewReader(s)
}
