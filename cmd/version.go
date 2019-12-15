// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

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
