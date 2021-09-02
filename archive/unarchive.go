package archive

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	// ErrNotSupportFile does not support file extensions.
	ErrNotSupportFile = errors.New("not support file")
)

// Unarchiver is a type that can extract archive files into a folder.
type Unarchiver interface {
	Unarchive(src, dest string) error
}

// Unarchive unarchives the given archive file into the destination folder.
// The archive format is selected implicitly.
func Unarchive(src, dest string) error {
	unarchiver, err := NewUnarchiver(src)
	if err != nil {
		return err
	}

	return unarchiver.Unarchive(src, dest)
}

// NewUnarchiver creates an unpacker that can extract archive files into a folder.
// The archive format is selected implicitly.
func NewUnarchiver(path string) (Unarchiver, error) {
	switch {
	case strings.HasSuffix(path, ".zip"):
		return Zip{}, nil

	case strings.HasSuffix(path, ".tar.gz"), strings.HasSuffix(path, ".tgz"):
		return TGz{}, nil

	default:
		return nil, fmt.Errorf("ext '%s': %w", filepath.Ext(path), ErrNotSupportFile)
	}
}

func invalidFilename(p string) bool {
	return p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../")
}
