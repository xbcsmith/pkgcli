// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

// Pkg struct for pkg
type Pkg struct {
	Description  string        `json:"description" yaml:"description"`
	Instructions []Instruction `json:"instructions" yaml:"instructions"`
	Name         string        `json:"name" yaml:"name"`
	Package      string        `json:"package" yaml:"package"`
	PlatformID   string        `json:"platform_id" yaml:"platform_id"`
	Provides     []string      `json:"provides" yaml:"provides"`
	Release      string        `json:"release" yaml:"release"`
	Requires     []string      `json:"requires" yaml:"requires"`
	Sources      []Source      `json:"sources" yaml:"sources"`
	Summary      string        `json:"summary" yaml:"summary"`
	Version      string        `json:"version" yaml:"version"`
}

// Pkgs struct for pkgs
type Pkgs struct {
	Packages []Pkg `json:"packages" yaml:"packages"`
}

// Instruction struct for instruction
type Instruction struct {
	Build     string `json:"build" yaml:"build"`
	Configure string `json:"configure" yaml:"configure"`
	Install   string `json:"install" yaml:"install"`
	Post      string `json:"post" yaml:"post"`
	Pre       string `json:"pre" yaml:"pre"`
	Test      string `json:"test" yaml:"test"`
	Unpack    string `json:"unpack" yaml:"unpack"`
}

// Instructions struct for instructions
type Instructions struct {
	Instructions []Instruction `json:"instructions" yaml:"instructions"`
}

// Source struct for source
type Source struct {
	Archive     string `json:"archive" yaml:"archive"`
	Destination string `json:"destination" yaml:"destination"`
	SHA256      string `json:"sha256" yaml:"sha256"`
	MD5         string `json:"md5" yaml:"md5"`
}

// Sources struct for sources
type Sources []struct {
	Sources []Source `json:"sources" yaml:"sources"`
}

// File struct for file
type File struct {
	Path   string `json:"path" yaml:"path"`
	Name   string `json:"name" yaml:"name"`
	Mode   string `json:"mode" yaml:"mode"`
	SHA256 string `json:"sha256" yaml:"sha256"`
}

// Files struct for files
type Files struct {
	Files []File `json:"files" yaml:"files"`
}
