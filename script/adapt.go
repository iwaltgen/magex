package script

import (
	"bytes"
	"io"
	"os"

	"github.com/bitfield/script"
)

// Pipe represents a pipe object with an associated ReadAutoCloser.
// https://pkg.go.dev/github.com/bitfield/script#readme-quick-start-unix-equivalents
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

// Args creates a pipe containing the program's command-line arguments, one per line.
var Args = script.Args

// IfExists tests whether the specified file exists, and creates a pipe whose
// error status reflects the result. If the file doesn't exist, the pipe's error
// status will be set, and if the file does exist, the pipe will have no error
// status. This can be used to do some operation only if a given file exists:
//
// IfExists("/foo/bar").Exec("/usr/bin/something")
var IfExists = script.IfExists

// Slice creates a pipe containing each element of the supplied slice of
// strings, one per line.
var Slice = script.Slice

// Buffer returns a *Pipe associated with the reader buffers.
// This is useful for starting pipelines.
func Buffer(buf io.Reader) *Pipe {
	return script.NewPipe().WithReader(buf)
}

// ReadFile returns a *Pipe associated with the reader file.
// This is useful for starting pipelines.
func ReadFile(name string) *Pipe {
	p := script.NewPipe()
	body, err := os.ReadFile(name)
	if err != nil {
		return p.WithError(err)
	}
	return p.WithReader(bytes.NewBuffer(body))
}
