package archive

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMkdir(t *testing.T) {
	dir, err := os.MkdirTemp("", "magex-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// dataset
	dataset := []string{"test1", "test2", "test3"}

	// table driven tests
	for _, v := range dataset {
		t.Run(v, func(t *testing.T) {
			// when
			err := mkdir(filepath.Join(dir, v), os.ModePerm)

			// then
			assert.NoError(t, err)
		})
	}
}

func TestWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "magex-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// dataset
	dataset := []string{"test1", "test2", "test3"}

	// table driven tests
	for _, v := range dataset {
		t.Run(v, func(t *testing.T) {
			// given
			path := filepath.Join(dir, v)
			body := bytes.NewBuffer([]byte("0123456789"))

			// when
			err := writeNewFile(path, body, os.ModePerm)

			// then
			assert.NoError(t, err)

			info, err := os.Lstat(path)
			assert.NoError(t, err)
			assert.NotNil(t, info)
		})
	}
}
