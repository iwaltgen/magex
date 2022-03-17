package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
)

// CreateTagOption describes how a tag object should be created.
type CreateTagOption func(*createTagOptions)

type createTagOptions struct {
	path    string                 // Default: "."
	refName plumbing.ReferenceName // Default: "HEAD"
	tag     *git.CreateTagOptions  // Default: message -> ""
	push    *git.PushOptions       // Issue: ssh-keyscan github.com >> ~/.ssh/known_hosts
	// TODO(iwaltgen): with update file option.
}

// CreateTag creates a tag. If opts is included, the tag is an annotated tag,
// otherwise a lightweight tag is created.
func CreateTag(tag string, opts ...CreateTagOption) error {
	opt := newCreateTagOptions(tag, opts...)
	repo, err := git.PlainOpen(opt.path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	ref, err := repo.Reference(opt.refName, true)
	if err != nil {
		return fmt.Errorf("failed to resolve reference: %w", err)
	}

	_, err = repo.CreateTag(tag, ref.Hash(), opt.tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	if opt.push != nil {
		if err := repo.Push(opt.push); err != nil {
			return fmt.Errorf("failed to push %s: %w", opt.push.RemoteName, err)
		}
	}
	return nil
}

// WithCreateTagPath set repository path.
func WithCreateTagPath(path string) CreateTagOption {
	return func(o *createTagOptions) {
		o.path = path
	}
}

// WithCreateTagRef adds a tag option reference.
func WithCreateTagRef(reference string) CreateTagOption {
	return func(o *createTagOptions) {
		o.refName = plumbing.ReferenceName(reference)
	}
}

// WithCreateTagMessage adds a tag option message.
func WithCreateTagMessage(message string) CreateTagOption {
	return func(o *createTagOptions) {
		o.tag.Message = message
	}
}

// WithCreateTagPush adds a push performs.
func WithCreateTagPush() CreateTagOption {
	return func(o *createTagOptions) {
		if o.push == nil {
			o.push = &git.PushOptions{}
		}
	}
}

// WithCreateTagRemote adds a push remote name.
func WithCreateTagRemote(remote string) CreateTagOption {
	return func(o *createTagOptions) {
		WithCreateTagPush()(o)
		o.push.RemoteName = remote
	}
}

// WithCreateTagRemote adds push progress.
func WithCreateTagProgress(progress sideband.Progress) CreateTagOption {
	return func(o *createTagOptions) {
		WithCreateTagPush()(o)
		o.push.Progress = progress
	}
}

func newCreateTagOptions(tag string, opts ...CreateTagOption) *createTagOptions {
	ret := &createTagOptions{
		path:    ".",
		refName: plumbing.HEAD,
		tag:     &git.CreateTagOptions{},
	}
	for _, v := range opts {
		v(ret)
	}

	if ret.push != nil {
		ret.push.RefSpecs = []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/tags/%[1]s:refs/tags/%[1]s", tag)),
		}
	}
	return ret
}
