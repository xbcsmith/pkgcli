package cmd

const (
	// AppName is application name
	AppName = "pkgcli"
	// AppVersion is the application version
	AppVersion = "dev"
	// AppRelease is the application release
	AppRelease = "dirty"
)

// GetVersion returns name-version-release as a string
func GetVersion() string {
	return AppName + "-" + AppVersion + "-" + AppRelease
}
