package archive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func mkdir(path string, mode os.FileMode) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir `%s`: %w", path, err)
	}
	return nil
}

func writeNewFile(path string, reader io.Reader, mode os.FileMode) error {
	dir := filepath.Dir(path)
	err := mkdir(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create file `%s`: %w", path, err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("write file `%s`: %w", path, err)
	}
	return nil
}
