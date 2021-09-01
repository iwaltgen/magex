package github

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/iwaltgen/magex/archive"
)

var defaultAPIGithubURL *url.URL

func init() {
	defaultAPIGithubURL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
	}
}

// Option represents download release file option.
type Option func(*option)

// DLReleaseFile downloads github content raw files.
func DLReleaseFile(repo string, opts ...Option) error {
	opt := newOption(opts...)
	url, err := latestReleaseFileURL(repo, opt.pattern)
	if err != nil {
		return err
	}

	filename := path.Base(url)
	target := filepath.Join(os.TempDir(), filename)
	if err := DownloadFile(url, target); err != nil {
		return err
	}
	defer os.Remove(target)

	if err := os.MkdirAll(opt.dest, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir '%s': %w", opt.dest, err)
	}

	switch {
	case opt.rename == "" && len(opt.pick) == 0:
		dest := filepath.Join(opt.dest, filename)
		if err := os.Rename(target, dest); err != nil {
			return fmt.Errorf("move '%s': %w", dest, err)
		}

	case opt.rename != "":
		rename := filepath.Join(opt.dest, opt.rename)
		if err := os.Rename(target, rename); err != nil {
			return fmt.Errorf("move '%s': %w", rename, err)
		}

	case len(opt.pick) != 0:
		tmp, err := os.MkdirTemp("", fmt.Sprintf("magex-*-%s", filename))
		if err != nil {
			return fmt.Errorf("mkdir tmp: %w", err)
		}
		defer os.RemoveAll(tmp)

		if err := archive.Unarchive(target, tmp); err != nil {
			return fmt.Errorf("unarchive '%s': %w", target, err)
		}

		if err := pickFiles(tmp, opt.dest, opt.pick); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unexpected case: %#v", opt)
	}
	return nil
}

// WithPattern adds the pattern to search assets option.
// Default pattern is `*{{.OS}}*{{.ARCH}}*`.
func WithPattern(pattern string) Option {
	return func(opt *option) {
		opt.pattern = pattern
	}
}

// WithDest represents downloads file local location option.
// Default dest is current directory('.').
func WithDest(dest string) Option {
	return func(opt *option) {
		opt.dest = dest
	}
}

// WithRename represents downloads file local rename option.
func WithRename(rename string) Option {
	return func(opt *option) {
		opt.rename = rename
	}
}

// WithPick represents downloads file inner select option.
func WithPick(files []string) Option {
	return func(opt *option) {
		opt.pick = map[string]string{}
		for _, v := range files {
			opt.pick[v] = v
		}
	}
}

// WithPickRename represents downloads file inner select and rename option.
func WithPickRename(pick map[string]string) Option {
	return func(opt *option) {
		opt.pick = pick
	}
}

type option struct {
	pattern string
	dest    string
	rename  string
	pick    map[string]string
}

func newOption(opts ...Option) *option {
	opt := &option{
		pattern: fmt.Sprintf("*%s*%s*", runtime.GOOS, runtime.GOARCH),
		dest:    ".",
	}

	for _, fn := range opts {
		fn(opt)
	}
	return opt
}

func latestReleaseFileURL(repo, pattern string) (string, error) {
	remote := *defaultAPIGithubURL
	remote.Path = path.Join("repos", repo, "releases/latest")

	url := remote.String()
	res, err := defaultClient.R().Get(url)
	if err != nil {
		return "", fmt.Errorf("request api '%s': %w", url, err)
	}

	selector := fmt.Sprintf(`assets.#(name%%"%s").browser_download_url`, pattern)
	value := gjson.GetBytes(res.Body(), selector)
	return value.String(), nil
}

func pickFiles(src, dest string, pick map[string]string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk '%s': %w", path, err)
		}
		if info.IsDir() {
			return nil
		}

		for origin, rename := range pick {
			if !strings.HasSuffix(path, origin) {
				continue
			}

			target := filepath.Join(dest, rename)
			if err := os.Rename(path, target); err != nil {
				return fmt.Errorf("move '%s' -> '%s': %w", path, target, err)
			}
		}
		return nil
	}
	return filepath.Walk(src, walkFn)
}
