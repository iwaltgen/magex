package semver

var defaultVersion = NewVersion("")

// Latest find latest version.
func Latest(versions []string) (string, error) {
	return defaultVersion.Latest(versions)
}

// Bump increase semantic version parts.
func Bump(version string, typ BumpType) (string, error) {
	return defaultVersion.Bump(version, typ)
}
