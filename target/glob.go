package target

import (
	"fmt"
	"os"

	"github.com/magefile/mage/target"
	"github.com/mattn/go-zglob"
)

// Glob expands each of the globs (file patterns) into individual sources and
// then calls Path on the result, reporting if any of the resulting sources have
// been modified more recently than the destination.
// Powered by glob: github.com/mattn/go-zglob
func Glob(dst string, globs ...string) (bool, error) {
	for _, g := range globs {
		files, err := zglob.Glob(os.ExpandEnv(g))
		if err != nil {
			return false, err
		}
		if len(files) == 0 {
			return false, fmt.Errorf("glob didn't match any files: %s", g)
		}

		shouldDo, err := target.Path(dst, files...)
		if err != nil {
			return false, err
		}
		if shouldDo {
			return true, nil
		}
	}
	return false, nil
}
