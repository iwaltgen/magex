package archive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func mkdir(path string, mode os.FileMode) error {
	err := os.MkdirAll(path, mode)
	if err != nil {
		return fmt.Errorf("mkdir `%s`: %w", path, err)
	}
	return nil
}

func writeNewFile(path string, reader io.Reader, mode os.FileMode) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir `%s`: %w", dir, err)
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
