package archive

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPickFiles(t *testing.T) {
	// dataset
	dataset := []struct {
		name  string
		files []string
		pick  map[string]string
	}{
		{
			name:  "single",
			files: []string{"file1", "file2"},
			pick: map[string]string{
				"file1": "file1",
			},
		},
		{
			name:  "single",
			files: []string{"file1", "inner/file2"},
			pick: map[string]string{
				"file1": "file",
			},
		},
		{
			name:  "single",
			files: []string{"file1", "inner/file2"},
			pick: map[string]string{
				"inner/file2": "file2",
			},
		},
		{
			name:  "single",
			files: []string{"file1", "inner/file2"},
			pick: map[string]string{
				"file2": "file2",
			},
		},
		{
			name:  "single",
			files: []string{"file1", "inner/file2"},
			pick: map[string]string{
				"file1": "inner/file1",
			},
		},
		{
			name:  "single",
			files: []string{"filename"},
			pick: map[string]string{
				"*": "file",
			},
		},
		{
			name:  "multiple",
			files: []string{"file1", "file2", "file3"},
			pick: map[string]string{
				"file1": "file1",
				"file2": "file2",
			},
		},
		{
			name:  "multiple",
			files: []string{"file1", "inner/file2", "in/in/file3"},
			pick: map[string]string{
				"file1": "file1",
				"file2": "file2",
			},
		},
		{
			name:  "multiple",
			files: []string{"file1", "inner/file2", "in/in/file3"},
			pick: map[string]string{
				"file2": "file2",
				"file3": "file3",
			},
		},
		{
			name:  "multiple",
			files: []string{"file1", "inner/file2", "in/in/file3"},
			pick: map[string]string{
				"file2": "file2",
				"file3": "inn/file3",
			},
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.name, func(t *testing.T) {
			// given
			src, err := os.MkdirTemp("", "magex-test-*")
			assert.NoError(t, err)
			dest, err := os.MkdirTemp("", "magex-test-*")
			assert.NoError(t, err)
			defer func() {
				os.RemoveAll(src)
				os.RemoveAll(dest)
			}()

			for _, v := range v.files {
				path := filepath.Join(src, v)
				body := bytes.NewBuffer([]byte(""))
				assert.NoError(t, writeNewFile(path, body, os.ModePerm))
			}

			// when
			err = PickFiles(src, dest, v.pick)

			// then
			assert.NoError(t, err)
			for k, v := range v.pick {
				if k != "*" {
					info, err := os.Lstat(filepath.Join(src, k))
					assert.True(t, os.IsNotExist(err))
					assert.Nil(t, info)
				}

				info, err := os.Lstat(filepath.Join(dest, v))
				assert.NoError(t, err)
				assert.NotNil(t, info)
			}
		})
	}
}

func TestPickFilesPermissionError(t *testing.T) {
	dest := "/home/unknown/magex"
	src, err := os.MkdirTemp("", "magex-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(src)

	files := []string{"file1", "file2"}
	for _, v := range files {
		path := filepath.Join(src, v)
		body := bytes.NewBuffer([]byte(""))
		assert.NoError(t, writeNewFile(path, body, os.ModePerm))
	}

	// when
	err = PickFiles(src, dest, map[string]string{
		"file1": "file1",
	})

	// then
	assert.Error(t, err)
}
