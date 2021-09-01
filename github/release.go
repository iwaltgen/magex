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
// Pattern: release name pattern. (default: '*{{.OS}}*{{.ARCH}}*')
func ReleaseAssetURL(repo, pattern string) (string, error) {
	if pattern == "" {
		pattern = fmt.Sprintf("*%s*%s*", runtime.GOOS, runtime.GOARCH)
	}

	return Release(repo, fmt.Sprintf(`assets.#(name%%"%s").browser_download_url`, pattern))
}

// PickReleaseAsset gets github latest release asset download url.
// Pattern: see ReleaseAssetURL
func PickReleaseAsset(repo, pattern string, opts ...http.Option) error {
	url, err := ReleaseAssetURL(repo, pattern)
	if err != nil {
		return err
	}

	return http.PickFile(url, opts...)
}
