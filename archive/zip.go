package archive

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

// Zip unarchives zip archive file.
type Zip struct{}

// Unarchive unpacks the .zip file from source to destination.
func (z Zip) Unarchive(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open reader '%s': %w", src, err)
	}
	defer r.Close()

	if err := mkdir(dest, os.ModePerm); err != nil {
		return err
	}

	for _, zf := range r.File {
		f, err := zf.Open()
		if err != nil {
			return fmt.Errorf("open file `%s`: %w", zf.Name, err)
		}
		defer f.Close()

		fileinfo := zf.FileInfo()
		fpath := filepath.Join(dest, zf.Name)
		if fileinfo.IsDir() {
			if err := mkdir(fpath, fileinfo.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := writeNewFile(fpath, f, fileinfo.Mode()); err != nil {
			return err
		}
	}
	return nil
}
