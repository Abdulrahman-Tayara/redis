package resp

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandReader_ReadCommand(t *testing.T) {
	command := "*3\r\n+set\r\n+key\r\n+hello\r\n"

	r := bytes.NewReader([]byte(command))

	cr := NewCommandReader(r, 4096)

	command, args, err := cr.ReadCommand()

	assert.Nil(t, err)

	assert.Equal(t, "set", command)

	assert.Equal(t, []any{"key", "hello"}, args)
}
