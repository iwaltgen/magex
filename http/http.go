package http

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"

	"github.com/iwaltgen/magex/archive"
)

var client *resty.Client

func init() {
	client = resty.New()
}

// Json is requests RESTful API then returns the response parsed value.
func Json(url, pattern string) (string, error) {
	res, err := client.R().
		SetHeader("accept", "application/json").
		Get(url)
	if err != nil {
		return "", err
	}

	value := gjson.GetBytes(res.Body(), pattern)
	return value.String(), nil
}

// Option represents unpack, pick files option.
type Option func(*option)

// File downloads file and unpack, pick files.
// Default: dest is local current directory.
func File(url string, opts ...Option) error {
	opt := newOption(opts...)
	filename := path.Base(url)
	target := filepath.Join(os.TempDir(), filename)
	if _, err := client.R().SetOutput(target).Get(url); err != nil {
		return err
	}
	defer os.Remove(target)

	if err := os.MkdirAll(opt.dir, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir '%s': %w", opt.dir, err)
	}

	switch {
	case opt.rename == "" && len(opt.pick) == 0:
		dest := filepath.Join(opt.dir, filename)
		if err := os.Rename(target, dest); err != nil {
			return fmt.Errorf("move '%s': %w", dest, err)
		}

	case opt.rename != "":
		rename := filepath.Join(opt.dir, opt.rename)
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

		if err := archive.PickFiles(tmp, opt.dir, opt.pick); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unexpected case: %#v", opt)
	}
	return nil
}

// WithDir represents downloads file local location option.
// Default dest is current directory('.').
func WithDir(dir string) Option {
	return func(opt *option) {
		opt.dir = dir
	}
}

// WithRename represents downloads file local rename option.
func WithRename(rename string) Option {
	return func(opt *option) {
		opt.rename = rename
	}
}

// WithPick represents downloads file inner pick option.
func WithPick(files []string) Option {
	return func(opt *option) {
		opt.pick = map[string]string{}
		for _, v := range files {
			opt.pick[v] = v
		}
	}
}

// WithPickRename represents downloads file inner pick and rename option.
func WithPickRename(pick map[string]string) Option {
	return func(opt *option) {
		opt.pick = pick
	}
}

type option struct {
	dir    string
	rename string
	pick   map[string]string
}

func newOption(opts ...Option) *option {
	opt := &option{
		dir: ".",
	}

	for _, fn := range opts {
		fn(opt)
	}
	return opt
}