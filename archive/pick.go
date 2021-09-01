package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PickFiles moves files from source to destination.
// Filename matches use 'HasSuffix'. '*' is all files.
func PickFiles(src, dest string, pick map[string]string) error {
	dest = os.ExpandEnv(dest)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk '%s': %w", path, err)
		}
		if info.IsDir() {
			return nil
		}

		for origin, rename := range pick {
			if origin != "*" && !strings.HasSuffix(path, origin) {
				continue
			}

			target := filepath.Join(dest, rename)
			if err := mkdir(filepath.Dir(target), os.ModePerm); err != nil {
				return err
			}
			if err := os.Rename(path, target); err != nil {
				return fmt.Errorf("move '%s' -> '%s': %w", path, target, err)
			}
			break
		}
		return nil
	}
	return filepath.Walk(os.ExpandEnv(src), walkFn)
}
