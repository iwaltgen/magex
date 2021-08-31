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

// RawFile represents remote/local file descriptor.
type RawFile struct {
	Remote string
	Local  string
}

// DLRawFile downloads github content raw files.
func DLRawFile(repo, branch string, files []RawFile) error {
	// TODO(iwaltgen): use multiple goroutine
	for _, f := range files {
		remote := *defaultGithubRawFileURL
		remote.Path = path.Join(repo, branch, f.Remote)
		dir := filepath.Dir(f.Local)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("mkdir(%s): %w", dir, err)
		}

		url := remote.String()
		if err := DownloadFile(url, f.Local); err != nil {
			return fmt.Errorf("download file(%s): %w", url, err)
		}
	}
	return nil
}
