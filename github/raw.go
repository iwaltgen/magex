package github

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

var defaultGithubRawFileURL *url.URL

func init() {
	defaultGithubRawFileURL = &url.URL{
		Scheme: "https",
		Host:   "raw.githubusercontent.com",
	}
}

// DLRawFile downloads github content raw files.
func DLRawFile(repo, branch string, files map[string]string) error {
	// TODO(iwaltgen): use multiple goroutine
	for remote, local := range files {
		url := *defaultGithubRawFileURL
		url.Path = path.Join(repo, branch, remote)
		dir := filepath.Dir(local)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("mkdir '%s': %w", dir, err)
		}

		if err := DownloadFile(url.String(), local); err != nil {
			return fmt.Errorf("download file '%v': %w", url, err)
		}
	}
	return nil
}
