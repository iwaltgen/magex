package github

import (
	"context"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/iwaltgen/magex/http"
	"golang.org/x/sync/errgroup"
)

var rawFileURL url.URL

func init() {
	rawFileURL = url.URL{
		Scheme: "https",
		Host:   "raw.githubusercontent.com",
	}
}

// RawFile downloads github content raw files.
func RawFile(repo, branch string, files map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	for remote, local := range files {
		url := rawFileURL
		url.Path = path.Join(repo, branch, remote)
		dest := os.ExpandEnv(local)

		g.Go(func() error {
			return http.File(url.String(), http.WithRename(dest))
		})
	}
	return g.Wait()
}
