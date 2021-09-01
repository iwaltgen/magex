package github

import (
	"fmt"
	"net/url"
	"path"
	"runtime"

	"github.com/iwaltgen/magex/http"
)

var apiURL *url.URL

func init() {
	apiURL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
	}
}

// Release gets github latest release info.
// Pattern: full API response select pattern.
func Release(repo, pattern string) (string, error) {
	remote := *apiURL
	remote.Path = path.Join("repos", repo, "releases/latest")

	url := remote.String()
	return http.Get(url, pattern)
}

// ReleaseAssetURL gets github latest release asset download url.
// Pattern: release name pattern.
func ReleaseAssetURL(repo string, opts ...Option) (string, error) {
	opt := &option{
		pattern: fmt.Sprintf("*%s*%s*", runtime.GOOS, runtime.GOARCH),
	}
	for _, fn := range opts {
		fn(opt)
	}

	return Release(repo, fmt.Sprintf(`assets.#(name%%"%s").browser_download_url`, opt.pattern))
}

// Option represents search assets option.
type Option func(*option)

// WithPattern adds the pattern to search assets.name option.
// Default pattern is `*{{.OS}}*{{.ARCH}}*`.
func WithPattern(pattern string) Option {
	return func(opt *option) {
		opt.pattern = pattern
	}
}

type option struct {
	pattern string
}
