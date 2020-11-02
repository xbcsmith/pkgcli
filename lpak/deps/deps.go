// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package deps

import "github.com/jinzhu/gorm"

// Dependencies is a representation of package deps
type Dependencies struct {
	gorm.Model
	Requires    []Dependency `json:"requires"`
	Provides    []Dependency `json:"provides"`
	Optional    []Dependency `json:"optional"`
	Recommended []Dependency `json:"recommended"`
}

// Dependency is a single entity that represents the package deps
type Dependency struct {
	gorm.Model
	PkgID     string `json:"pkgid"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Installed bool   `json:"installed"`
}
