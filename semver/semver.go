package semver

var defaultVersion = NewVersion("v")

// Latest find latest git tag version.
func LatestTag(path string) (string, error) {
	return defaultVersion.LatestTag(path)
}

// Latest find latest version.
func Latest(versions []string) (string, error) {
	return defaultVersion.Latest(versions)
}

// Bump increase semantic version parts.
func Bump(version string, typ BumpType) (string, error) {
	return defaultVersion.Bump(version, typ)
}
