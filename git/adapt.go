package git

import (
	"github.com/go-git/go-git/v5"
)

// Worktree represents a git worktree.
type Worktree = git.Worktree

// CheckoutOptions describes how a checkout operation should be performed.
type CheckoutOptions = git.CheckoutOptions

// CommitOptions describes how a commit operation should be performed.
type CommitOptions = git.CommitOptions

// PullOptions describes how a pull should be performed.
type PullOptions = git.PullOptions

// PushOptions describes how a push should be performed.
type PushOptions = git.PushOptions
