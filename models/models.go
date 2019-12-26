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
	Optional     []string      `json:"optional,omitempty" yaml:"optional,omitempty"`
	Recommended  []string      `json:"recommended,omitempty" yaml:"recommended,omitempty"`
	Sources      []Source      `json:"sources" yaml:"sources"`
	Files        []File        `json:"files" yaml:"files"`
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
	Configure string `json:"configure,omitempty" yaml:"configure,omitempty"`
	Install   string `json:"install" yaml:"install"`
	Post      string `json:"post,omitempty" yaml:"post,omitempty"`
	Pre       string `json:"pre,omitempty" yaml:"pre,omitempty"`
	Test      string `json:"test,omitempty" yaml:"test,omitempty"`
	Unpack    string `json:"unpack" yaml:"unpack"`
}

// Instructions struct for instructions
type Instructions struct {
	Instructions []Instruction `json:"instructions" yaml:"instructions"`
}

// Source struct for source
type Source struct {
	Archive     string `json:"archive" yaml:"archive"`
	Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`
	SHA256      string `json:"sha256,omitempty" yaml:"sha256,omitempty"`
	MD5         string `json:"md5,omitempty" yaml:"md5,omitempty"`
	OnDisk      string `json:"ondisk,omitempty" yaml:"ondisk,omitempty"`
	Size        string `json:"size,omitempty" yaml:"size,omitempty"`
}

// Sources struct for sources
type Sources []struct {
	Sources []Source `json:"sources" yaml:"sources"`
}

// File struct for file
type File struct {
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Mode   string `json:"mode,omitempty" yaml:"mode,omitempty"`
	SHA256 string `json:"sha256,omitempty" yaml:"sha256,omitempty"`
}

// Files struct for files
type Files struct {
	Files []File `json:"files" yaml:"files"`
}
