package github

import (
	"fmt"
	"net/url"
	"path"

	"github.com/iwaltgen/magex/http"
)

var rawFileURL *url.URL

func init() {
	rawFileURL = &url.URL{
		Scheme: "https",
		Host:   "raw.githubusercontent.com",
	}
}

// RawFile downloads github content raw files.
func RawFile(repo, branch string, files map[string]string) error {
	// TODO(iwaltgen): use multiple goroutine
	for remote, local := range files {
		url := *rawFileURL
		url.Path = path.Join(repo, branch, remote)
		if err := http.File(url.String(), http.WithRename(local)); err != nil {
			return fmt.Errorf("download file '%v': %w", url, err)
		}
	}
	return nil
}
