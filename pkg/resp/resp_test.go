package resp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshal_MarshalString(t *testing.T) {
	tests := []marshalTest[string]{
		{
			"Valid simple string",
			"OK",
			"+OK\r\n",
			nil,
		},
		{
			"Valid simple string with whitespaces",
			"OK OK2 ",
			"+OK OK2 \r\n",
			nil,
		},
		{
			"Valid bulk strings",
			"This is a bulk strings that\ncontains newlines and\r",
			"$50\r\nThis is a bulk strings that\ncontains newlines and\r\r\n",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, marshalTestRunner(&test))
	}
}

func TestMarshal_MarshalError(t *testing.T) {
	tests := []marshalTest[error]{
		{
			"Valid simple error",
			errors.New("some critical error"),
			"-some critical error\r\n",
			nil,
		},
		{
			"Invalid simple error - contains \\n",
			errors.New("some critical error\ndetails here..."),
			"",
			errors.New("simple string mustn't contain \\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, marshalTestRunner(&test))
	}
}

func TestMarshal_MarshalInt(t *testing.T) {
	tests := []marshalTest[any]{
		{
			"Valid int",
			10,
			":10\r\n",
			nil,
		},
		{
			"Valid negative int",
			-19,
			":-19\r\n",
			nil,
		},
		{
			"Valid int64",
			int64(10),
			":10\r\n",
			nil,
		},
		{
			"Valid negative int64",
			int64(-19),
			":-19\r\n",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, marshalTestRunner(&test))
	}
}

func TestMarshal_MarshalArray(t *testing.T) {
	tests := []marshalTest[any]{
		{
			"Valid int array",
			[]int{1, 2, 3},
			"*3\r\n:1\r\n:2\r\n:3\r\n",
			nil,
		},
		{
			"Valid string array",
			[]string{"Hello", "This is \r\nBulk strings", "World"},
			"*3\r\n+Hello\r\n$22\r\nThis is \r\nBulk strings\r\n+World\r\n",
			nil,
		},
		{
			"Valid empty array",
			[]any{},
			"*0\r\n",
			nil,
		},
		{
			"Valid mixed array",
			[]any{1, 3, "OK"},
			"*3\r\n:1\r\n:3\r\n+OK\r\n",
			nil,
		},
		{
			"Valid nested array",
			[]any{1, 3, []any{1, 3, "OK"}},
			"*3\r\n:1\r\n:3\r\n*3\r\n:1\r\n:3\r\n+OK\r\n",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, marshalTestRunner(&test))
	}
}

func TestMarshal_MarshalNull(t *testing.T) {
	tests := []marshalTest[any]{
		{
			"Valid null",
			nil,
			"_\r\n",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, marshalTestRunner(&test))
	}
}

type marshalTest[T any] struct {
	name        string
	input       T
	expected    string
	expectedErr error
}

func marshalTestRunner[T any](test *marshalTest[T]) func(t *testing.T) {
	return func(t *testing.T) {
		output, err := Marshal(test.input)

		assert.Equal(t, test.expectedErr, err)

		if err == nil {
			assert.Equal(t, test.expected, output)
		}
	}
}
