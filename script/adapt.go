package script

import (
	"io"

	"github.com/bitfield/script"
)

// Pipe represents a pipe object with an associated ReadAutoCloser.
type Pipe = script.Pipe

// File returns a *Pipe associated with the specified file. This is useful for
// starting pipelines. If there is an error opening the file, the pipe's error
// status will be set.
var File = script.File

// FindFiles takes a directory path and returns a pipe listing all the files in
// the directory and its subdirectories recursively, one per line, like Unix
// `find -type f`. If the path doesn't exist or can't be read, the pipe's error
// status will be set.
var FindFiles = script.FindFiles

// ListFiles creates a pipe containing the files and directories matching the
// supplied path, one per line. The path may be a glob, conforming to
// filepath.Match syntax.
var ListFiles = script.ListFiles

// Echo returns a pipe containing the supplied string.
var Echo = script.Echo

// Buffer returns a *Pipe associated with the reader buffers. This is useful for
// starting pipelines. If there is an error opening the file, the pipe's error
// status will be set.
func Buffer(buf io.Reader) *Pipe {
	return script.NewPipe().WithReader(buf)
}
