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
