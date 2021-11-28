package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawFile(t *testing.T) {
	dir := "testdata"
	err := os.Mkdir(dir, os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// dataset
	dataset := []struct {
		name         string
		repo, branch string
		files        map[string]string
	}{
		{
			name:   "envoyproxy/pgv",
			repo:   "envoyproxy/protoc-gen-validate",
			branch: "master",
			files: map[string]string{
				"validate/validate.proto": filepath.Join(dir, "api/envoyproxy/pgv/validate.proto"),
			},
		},
		{
			name:   "googleapis/rpc",
			repo:   "googleapis/googleapis",
			branch: "master",
			files: map[string]string{
				"google/rpc/code.proto":          filepath.Join(dir, "api/google/rpc/code.proto"),
				"google/rpc/error_details.proto": filepath.Join(dir, "api/google/rpc/error_details.proto"),
				"google/rpc/status.proto":        filepath.Join(dir, "api/google/rpc/status.proto"),
			},
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.name, func(t *testing.T) {
			// when
			err := RawFile(v.repo, v.branch, v.files)

			// then
			assert.NoError(t, err)
		})
	}
}

func TestRawFileNotFound(t *testing.T) {
	dir := "testdata"
	err := os.Mkdir(dir, os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// when
	repo := "envoyproxy/protoc-gen-validate"
	branch := "main"
	files := map[string]string{
		"validate/validate.not_exists.proto": filepath.Join(dir, "api/envoyproxy/pgv/validate.proto"),
	}
	err = RawFile(repo, branch, files)

	// then
	assert.Error(t, err)
}
