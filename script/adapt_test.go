package script

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	// when
	pipe := Buffer(bytes.NewBuffer([]byte("test1\ntest2")))
	list, err := pipe.Slice()

	// then
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
