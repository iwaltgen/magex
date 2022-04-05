//go:generate stringer -type=BumpType -output=bump_type_string.go

package semver

import "strings"

// BumpType increases the version section number type.
type BumpType int32

const (
	Unspecific BumpType = iota
	Major
	Minor
	Patch
)

// ParseBumpType parse string bump type.
func ParseBumpType(str string) BumpType {
	switch strings.ToLower(str) {
	case "major":
		return Major

	case "minor":
		return Minor

	case "patch":
		return Patch

	default:
		return Unspecific
	}
}
