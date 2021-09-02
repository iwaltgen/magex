package dep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadImport(t *testing.T) {
	// dataset
	dataset := []struct {
		name     string
		patterns []string
		paths    []string
	}{
		{
			name:     "single",
			patterns: []string{"testdata/deps1.go"},
			paths: []string{
				"github.com/golangci/golangci-lint/cmd/golangci-lint",
				"golang.org/x/tools/cmd/stringer",
			},
		},
		{
			name:     "single",
			patterns: []string{"testdata/deps2.go"},
			paths: []string{
				"github.com/mfridman/tparse",
				"golang.org/x/tools/cmd/stringer",
			},
		},
		{
			name:     "multiple",
			patterns: []string{"testdata/deps1.go", "testdata/deps2.go"},
			paths: []string{
				"github.com/golangci/golangci-lint/cmd/golangci-lint",
				"github.com/mfridman/tparse",
				"golang.org/x/tools/cmd/stringer",
			},
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.name, func(t *testing.T) {
			// when
			paths, err := GlobImport(v.patterns...)

			// then
			assert.NoError(t, err)
			assert.Equal(t, v.paths, paths)
		})
	}
}

func TestLoadImportNotExist(t *testing.T) {
	// when
	ret, err := GlobImport("testdata/test1.go")

	// then
	assert.Error(t, err)
	assert.Nil(t, ret)
}

func TestLoadImportNotMatch(t *testing.T) {
	// when
	ret, err := GlobImport("testdata/test*")

	// then
	assert.Error(t, err)
	assert.Nil(t, ret)
}
