package resp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshaller_UnmarshalString(t *testing.T) {
	tests := []unmarshalTest[string]{
		{
			"Valid simple string",
			"+OK\r\n",
			"OK",
			nil,
		},
		{
			"Invalid simple string - missing CRLF",
			"+OK",
			"",
			ErrValueNotEndWithCrlf,
		},
		{
			"Invalid simple string - contains LF",
			"+OK\nOK2\r\n",
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
			"-Error message\r\n",
			"Error message",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

func TestUnmarshaller_UnmarshalInt(t *testing.T) {
	tests := []unmarshalTest[int]{
		{
			"Valid int",
			":10\r\n",
			10,
			nil,
		},
		{
			"Valid signed int",
			":-12\r\n",
			-12,
			nil,
		},
		{
			"Invalid int",
			":-12E\r\n",
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
			"$16\r\nHello\rMy\nFriends\r\n",
			"Hello\rMy\nFriends",
			nil,
		},
		{
			"Valid empty bulk strings",
			"$0\r\n\r\n",
			"",
			nil,
		},
		{
			"Invalid bulk strings - invalid string length",
			"$Hello\r\nWorld\r\n",
			"",
			errors.New("invalid string length value, expected one or more decimal digits (0..9)"),
		},
		{
			"Invalid bulk strings - mismatch content length with passed string length",
			"$2\r\nWorld\r\n",
			"",
			errors.New("content length doesn't match the passed string length"),
		},
		{
			"Invalid bulk strings - invalid bulk strings structure",
			"$2\r\nHello\r\nWorld\r\n",
			"",
			errors.New("invalid bulk strings structure"),
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
			"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			[]any{"hello", "world"},
			nil,
		},
		{
			"Valid empty array",
			"*0\r\n",
			[]any{},
			nil,
		},
		{
			"Valid multi-type array",
			"*4\r\n$5\r\nhello\r\n:12\r\n-Error Message\r\n+OK\r\n",
			[]any{"hello", 12, "Error Message", "OK"},
			nil,
		},
		{
			"Valid nested arrays",
			"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n",
			[]any{[]any{1, 2, 3}, []any{"Hello", "World"}},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, unmarshalTestRunner(&test))
	}
}

type unmarshalTest[T any] struct {
	name        string
	input       any
	expected    T
	expectedErr error
}

func unmarshalTestRunner[T any](test *unmarshalTest[T]) func(t *testing.T) {
	return func(t *testing.T) {
		output, err := Unmarshal(test.input)

		assert.Equal(t, test.expectedErr, err)

		if err == nil {
			assert.Equal(t, test.expected, output)
		}
	}
}

func TestParseArrayElements(t *testing.T) {
	parts := []string{"*2\r\n", "*3\r\n", ":1\r\n", ":2\r\n", ":3\r\n", "*2\r\n", "+Hello\r\n", "-World\r\n"}

	expected := []string{"*3\r\n:1\r\n:2\r\n:3\r\n", "*2\r\n+Hello\r\n-World\r\n"}

	output, err := parseArrayElements(parts)

	assert.NoError(t, err)

	assert.Equal(t, expected, output)
}
