package target

import "github.com/magefile/mage/target"

// Path first expands environment variables like $FOO or ${FOO}, and then
// reports if any of the sources have been modified more recently than the
// destination. Path does not descend into directories, it literally just checks
// the modtime of each thing you pass to it. If the destination file doesn't
// exist, it always returns true and nil. It's an error if any of the sources
// don't exist.
var Path = target.Path

// Dir reports whether any of the sources have been modified more recently
// than the destination. If a source or destination is a directory, this
// function returns true if a source has any file that has been modified more
// recently than the most recently modified file in dst. If the destination
// file doesn't exist, it always returns true and nil.  It's an error if any
// of the sources don't exist.
var Dir = target.Dir
