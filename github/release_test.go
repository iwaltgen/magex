package github

import (
	"os"
	"testing"

	"github.com/iwaltgen/magex/http"
	"github.com/stretchr/testify/assert"
)

func TestRelease(t *testing.T) {
	// dataset
	dataset := []struct {
		repo    string
		pattern string
	}{
		{
			repo:    "hashicorp/consul",
			pattern: PatternName,
		},
		{
			repo:    "prometheus/prometheus",
			pattern: PatternTagName,
		},
		{
			repo:    "prometheus/alertmanager",
			pattern: PatternAssetCurrentOS,
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.repo, func(t *testing.T) {
			// when
			ret, err := Release(v.repo, v.pattern)

			// then
			assert.NoError(t, err)
			assert.NotEmpty(t, ret)
		})
	}
}

func TestReleaseNotFound(t *testing.T) {
	// dataset
	dataset := []struct {
		repo    string
		pattern string
	}{
		{
			repo:    "googleapis/googleapis",
			pattern: PatternID,
		},
		{
			repo:    "iwaltgen/magex",
			pattern: PatternURL,
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.repo, func(t *testing.T) {
			// when
			ret, err := Release(v.repo, v.pattern)

			// then
			assert.Error(t, err)
			assert.Empty(t, ret)
		})
	}
}

func TestReleaseAsset(t *testing.T) {
	dir, err := os.MkdirTemp("", "magex-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// dataset
	dataset := []struct {
		repo    string
		pattern string
		opts    []http.Option
	}{
		{
			repo:    "FiloSottile/age",
			pattern: PatternAssetCurrentOS,
		},
		{
			repo:    "FiloSottile/age",
			pattern: PatternAssetCurrentOS,
			opts: []http.Option{
				http.WithPick("age", "age-keygen"),
			},
		},
		{
			repo:    "FiloSottile/age",
			pattern: PatternAssetCurrentOS,
			opts: []http.Option{
				http.WithPickRename(map[string]string{
					"age":        "age",
					"age-keygen": "keygen",
				}),
			},
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.repo, func(t *testing.T) {
			// when
			opts := []http.Option{http.WithDir(dir)}
			opts = append(opts, v.opts...)
			err := ReleaseAsset(v.repo, v.pattern, opts...)

			// then
			assert.NoError(t, err)
		})
	}
}

func TestReleaseAssetNotFound(t *testing.T) {
	// dataset
	dataset := []struct {
		repo    string
		pattern string
	}{
		{
			repo:    "googleapis/googleapis",
			pattern: PatternID,
		},
		{
			repo:    "iwaltgen/magex",
			pattern: PatternURL,
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.repo, func(t *testing.T) {
			// when
			err := ReleaseAsset(v.repo, v.pattern)

			// then
			assert.Error(t, err)
		})
	}
}

func TestReleaseAssetDirPermissionError(t *testing.T) {
	// dataset
	dataset := []struct {
		repo    string
		pattern string
		opts    []http.Option
	}{
		{
			repo:    "FiloSottile/age",
			pattern: PatternAssetCurrentOS,
			opts: []http.Option{
				http.WithDir("/home/unknown/magex"),
			},
		},
		{
			repo:    "FiloSottile/age",
			pattern: PatternAssetCurrentOS,
			opts: []http.Option{
				http.WithDir("/home/unknown/magex"),
				http.WithRename("age.tar.gz"),
			},
		},
	}

	// table driven tests
	for _, v := range dataset {
		t.Run(v.repo, func(t *testing.T) {
			// when
			err := ReleaseAsset(v.repo, v.pattern, v.opts...)

			// then
			assert.Error(t, err)
		})
	}
}
