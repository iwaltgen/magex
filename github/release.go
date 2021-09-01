package github

import (
	"fmt"
	"net/url"
	"path"
	"runtime"

	"github.com/iwaltgen/magex/http"
)

var apiURL url.URL

func init() {
	apiURL = url.URL{
		Scheme: "https",
		Host:   "api.github.com",
	}
}

// https://pkg.go.dev/github.com/tidwall/gjson#readme-path-syntax
// https://docs.github.com/en/rest/reference/repos#get-the-latest-release
const (
	PatternID             = `id`
	PatternName           = `name`
	PatternURL            = `url`
	PatternHtmlURL        = `html_url`
	PatternTagName        = `tag_name`
	PatternBody           = `body`
	PatternAssetCurrentOS = `assets.#(name%%"` + "*" + runtime.GOOS + "*" + runtime.GOARCH + "*" + `").browser_download_url`
)

// Release gets github latest release info.
// pattern: PatternID, PatternTagName, PatternAssetCurrentOS...
func Release(repo, pattern string) (string, error) {
	url := apiURL
	url.Path = path.Join("repos", repo, "releases/latest")
	return http.Json(url.String(), pattern)
}

// ReleaseAsset gets github latest release asset download url.
// pattern: PatternAssetCurrentOS...
func ReleaseAsset(repo, pattern string, opts ...http.Option) error {
	url, err := Release(repo, pattern)
	if err != nil {
		return err
	}

	return http.File(url, opts...)
}

// PatternAssetURL makes releases asset name pattern for download URL.
func PatternAssetURL(pattern string) string {
	return fmt.Sprintf(`assets.#(name%%"%s").browser_download_url`, pattern)
}
