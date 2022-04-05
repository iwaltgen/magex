package semver

import (
	"fmt"

	"github.com/Masterminds/semver"
)

// Version implements semantic versioning handle.
type Version struct {
	prefix string
}

// NewVersion creates version handle.
func NewVersion(prefix string) Version {
	return Version{
		prefix: prefix,
	}
}

// Latest find latest version.
func (v Version) Latest(versions []string) (string, error) {
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
	return v.prefix + ret.String(), nil
}

// Bump increase semantic version parts.
func (v Version) Bump(version string, typ BumpType) (string, error) {
	current, _ := semver.NewVersion(version)
	var next semver.Version

	switch typ {
	case Major:
		next = current.IncMajor()

	case Minor:
		next = current.IncMinor()

	case Patch:
		next = current.IncPatch()

	default:
		return "", errTypeWithUsage(typ)
	}

	return v.prefix + next.String(), nil
}

func errTypeWithUsage(typ BumpType) error {
	return fmt.Errorf(`invalid bump version type: %s

Semantic Versioning (https://semver.org)
major: bump up next major version
minor: bump up next minor version
patch: bump up next patch version
`, typ)
}
