package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	// dataset
	dataset := []struct {
		name string
		url  string
	}{
		{
			name: "mage",
			url:  "https://api.github.com/repos/magefile/mage",
		},
		{
			name: "magex",
			url:  "https://api.github.com/repos/iwaltgen/magex",
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.name, func(t *testing.T) {
			// when
			name, err := Json(v.url, "name")

			// then
			assert.NoError(t, err)
			assert.NotEmpty(t, name)
		})
	}
}

func TestJsonNotFound(t *testing.T) {
	// when
	url := "https://api.github.com/repos/iwaltgen/magex/releases/latest"
	version, err := Json(url, "tag_name")

	// then
	assert.Error(t, err)
	assert.Empty(t, version)
}
