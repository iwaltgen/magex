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
			name:     "Single",
			patterns: []string{"testdata1.go"},
			paths: []string{
				"github.com/golangci/golangci-lint/cmd/golangci-lint",
				"golang.org/x/tools/cmd/stringer",
			},
		},
		{
			name:     "Single",
			patterns: []string{"testdata2.go"},
			paths: []string{
				"github.com/mfridman/tparse",
				"golang.org/x/tools/cmd/stringer",
			},
		},
		{
			name:     "Multiple",
			patterns: []string{"testdata1.go", "testdata2.go"},
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
