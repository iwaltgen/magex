package http

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "magex-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// dataset
	dataset := []struct {
		name string
		url  string
		opts []Option
	}{
		{
			name: "tarball",
			url:  "https://api.github.com/repos/magefile/mage/tarball/v1.11.0",
		},
		{
			name: "tarball",
			url:  "https://api.github.com/repos/magefile/mage/tarball/v1.11.0",
			opts: []Option{WithRename("tarball-rename.tar.gz")},
		},
		{
			name: "tarball",
			url:  "https://api.github.com/repos/magefile/mage/tarball/v1.11.0",
			opts: []Option{WithRename("tarball.tar.gz"), WithPick("picktar/go.mod", "picktar/go.sum")},
		},
		{
			name: "zipball",
			url:  "https://api.github.com/repos/magefile/mage/zipball/v1.11.0",
		},
		{
			name: "zipball",
			url:  "https://api.github.com/repos/magefile/mage/zipball/v1.11.0",
			opts: []Option{WithRename("zipball-rename.zip")},
		},
		{
			name: "zipball",
			url:  "https://api.github.com/repos/magefile/mage/zipball/v1.11.0",
			opts: []Option{WithRename("zipball.zip"), WithPickRename(map[string]string{
				"magefile.go": "pickzip/mmagefile.go",
				"LICENSE":     "pickzip/LICENSE_MAGE",
			})},
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.name, func(t *testing.T) {
			// when
			opts := []Option{WithDir(dir)}
			opts = append(opts, v.opts...)
			err := File(v.url, opts...)

			// then
			assert.NoError(t, err)
		})
	}
}

func TestFileNotFound(t *testing.T) {
	// when
	url := "https://api.github.com/repos/magefile/mage/tarball/v1.9.10"
	err := File(url)

	// then
	assert.Error(t, err)
}
