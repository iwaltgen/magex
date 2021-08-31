package semver

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
)

// Bump version
func Bump(version, typ string) (string, error) {
	current, _ := semver.NewVersion(version)
	var next semver.Version

	switch strings.ToLower(typ) {
	case "major":
		next = current.IncMajor()

	case "minor":
		next = current.IncMinor()

	case "patch":
		next = current.IncPatch()

	default:
		return "", errTypeWithUsage(typ)
	}

	return next.String(), nil
}

func errTypeWithUsage(typ string) error {
	return fmt.Errorf(`invalid bump version type: %s

Semantic Versioning (https://semver.org)
major: bump up next major version
minor: bump up next minor version
patch: bump up next patch version
`, typ)
}
