package semver

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
)

// Latest find latest version.
func Latest(versions []string) (string, error) {
	var ret *semver.Version
	for _, v := range versions {
		ver, err := semver.NewVersion(v)
		if err != nil {
			return "", err
		}

		if ret == nil {
			ret = ver
			continue
		}

		if ver.GreaterThan(ret) {
			ret = ver
		}
	}
	return ret.String(), nil
}

// Bump increase semantic version parts.
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
