package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TGz unarchives tar.gz(tgz) archive file.
type TGz struct{}

// Unarchive unpacks the .tar.gz(tgz) file from source to destination.
func (t TGz) Unarchive(src, dest string) error {
	src = os.ExpandEnv(src)
	dest = os.ExpandEnv(dest)

	sf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file '%s': %w", src, err)
	}
	defer sf.Close()

	gr, err := gzip.NewReader(sf)
	if err != nil {
		return fmt.Errorf("open reader: %w", err)
	}
	defer gr.Close()

	if err := mkdir(dest, os.ModePerm); err != nil {
		return err
	}

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reader next: %w", err)
		}

		if invalidFilename(header.Name) {
			return fmt.Errorf("invalid file name '%s'", header.Name)
		}

		path := filepath.Join(dest, header.Name)
		mode := os.FileMode(header.Mode)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := mkdir(path, mode); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := writeNewFile(path, tr, mode); err != nil {
				return err
			}
		}
	}
	return nil
}
