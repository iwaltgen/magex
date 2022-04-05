package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
)

// Repository represents a git repository
type Repository = git.Repository

// Worktree represents a git worktree.
type Worktree = git.Worktree

// NewRepository opens a git repository from the given path. It detects if the
// repository is bare or a normal one. If the path doesn't contain a valid
// repository ErrRepositoryNotExists is returned
func NewRepository(path string) (*Repository, error) {
	return git.PlainOpen(path)
}

// Tags returns all the tag References in a repository.
func Tags(path string) ([]string, error) {
	repo, err := NewRepository(path)
	if err != nil {
		return nil, err
	}

	refs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var ret []string
	err = refs.ForEach(func(r *plumbing.Reference) error {
		ret = append(ret, r.Name().Short())
		return nil
	})
	return ret, err
}

// CreateTagOption describes how a tag object should be created.
type CreateTagOption func(*createTagOptions)

type createTagOptions struct {
	path    string                 // Default: "."
	refName plumbing.ReferenceName // Default: "HEAD"
	tag     *git.CreateTagOptions  // Default: message -> ""
	push    *git.PushOptions       // Issue: ssh-keyscan github.com >> ~/.ssh/known_hosts
	hook    func(*Repository) error
}

// CreateTag creates a tag. If opts is included, the tag is an annotated tag,
// otherwise a lightweight tag is created.
func CreateTag(tag string, opts ...CreateTagOption) error {
	opt := newCreateTagOptions(tag, opts...)
	repo, err := NewRepository(opt.path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	if err := opt.hook(repo); err != nil {
		return fmt.Errorf("failed to hook: %w", err)
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

// WithCreateTagHook adds a tag create hook apply files.
func WithCreateTagHook(fn func(*Repository) error) CreateTagOption {
	return func(o *createTagOptions) {
		o.hook = fn
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

// WithCreateTagProgress adds push progress.
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
		hook: func(r *Repository) error {
			return nil
		},
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
