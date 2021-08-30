package dep

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
	"golang.org/x/tools/go/packages"
)

// GlobImport expands each of the globs (file patterns) into individual sources
// and gets an import library path.
// Powered by glob: github.com/mattn/go-zglob
func GlobImport(globs ...string) ([]string, error) {
	imports := map[string]struct{}{}
	for _, g := range globs {
		files, err := zglob.Glob(g)
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("failed to glob didn't match any files: %s", g)
		}

		ret, err := loadImport(files...)
		if err != nil {
			return nil, fmt.Errorf("failed to load import(%s): %w", g, err)
		}

		for _, v := range ret {
			imports[v] = struct{}{}
		}
	}

	ret := make([]string, 0, len(imports))
	for path := range imports {
		ret = append(ret, path)
	}
	sort.Strings(ret)
	return ret, nil
}

func loadImport(files ...string) ([]string, error) {
	pkgs, err := packages.Load(config(), files...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	imports := []string{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(node ast.Node) bool {
				file, ok := node.(*ast.File)
				if !ok {
					return true
				}

				for _, v := range file.Imports {
					imports = append(imports, strings.Trim(v.Path.Value, "\""))
				}
				return false
			})
		}
	}

	return imports, nil
}

func config() *packages.Config {
	mode := packages.NeedFiles
	mode |= packages.NeedImports
	mode |= packages.NeedSyntax
	return &packages.Config{
		Mode:  mode,
		Tests: false,
	}
}
