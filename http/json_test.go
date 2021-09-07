package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestJson(t *testing.T) {
	// when
	ret, err := Json(
		"https://api.github.com/repos/magefile/mage/releases/latest",
		"assets.#.name",
	)

	// then
	assert.NoError(t, err)
	assert.Equal(t, gjson.JSON, ret.Type)
	for _, v := range ret.Array() {
		assert.Contains(t, v.String(), "mage")
	}
}

func TestJsonString(t *testing.T) {
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
			name, err := JsonString(v.url, "name")

			// then
			assert.NoError(t, err)
			assert.NotEmpty(t, name)
		})
	}
}

func TestJsonNotFound(t *testing.T) {
	// when
	url := "https://api.github.com/repos/iwaltgen/magex/releases/latest"
	tag, err := JsonString(url, "tag_name")

	// then
	assert.Error(t, err)
	assert.Empty(t, tag)
}
