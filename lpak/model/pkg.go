// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"crypto/md5" // nolint:gosec
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"

	"github.com/xbcsmith/pkgcli/lpak/common"
	"github.com/xbcsmith/pkgcli/lpak/deps"
	"github.com/xbcsmith/pkgcli/lpak/files"
	"github.com/xbcsmith/pkgcli/lpak/instructions"
	"github.com/xbcsmith/pkgcli/lpak/source"
	yaml "gopkg.in/yaml.v3"
)

// Pkg struct for pkg
type Pkg struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	Updated      int64  `gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
	Created      int64  `gorm:"autoCreateTime"`      // Use unix seconds as creating time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time                  `gorm:"index"`
	Installed    bool                       `json:"installed" yaml:"installed"`
	PkgID        string                     `json:"pkgid" yaml:"pkgid"`
	Description  string                     `json:"description" yaml:"description"`
	Name         string                     `json:"name" yaml:"name"`
	Version      string                     `json:"version" yaml:"version"`
	Package      string                     `json:"package" yaml:"package"`
	Arch         string                     `json:"arch" yaml:"arch"`
	Platform     string                     `json:"platform" yaml:"platform"`
	Summary      string                     `json:"summary" yaml:"summary"`
	Release      string                     `json:"release" yaml:"release"`
	Provides     []deps.Dependency          `json:"provides" yaml:"provides"`
	Requires     []deps.Dependency          `json:"requires" yaml:"requires"`
	Optional     []deps.Dependency          `json:"optional,omitempty" yaml:"optional,omitempty"`
	Recommended  []deps.Dependency          `json:"recommended,omitempty" yaml:"recommended,omitempty"`
	Instructions []instructions.Instruction `json:"instructions" yaml:"instructions"`
	Sources      []source.Source            `json:"sources" yaml:"sources"`
	Files        []files.File               `json:"files" yaml:"files"`
}

// Pkgs struct for pkgs
type Pkgs struct {
	Packages []Pkg `json:"packages" yaml:"packages"`
}

// NewPkg func takes name and version as input and returns *Pkg
func NewPkg(name, version string) *Pkg {
	release := common.NewRelease()
	description := strings.Title(name) + " " + version + " " + release
	return &Pkg{
		Name:        name,
		Version:     version,
		Release:     release,
		Description: description,
		Summary:     description,
		Package:     "tar.xz",
		Platform:    "x86_64-gnu-linux-9",
		PkgID:       common.NewULIDAsString(),
	}
}

// GetNVRA returns the Name Version Release Arch of a package.
// foo-1.0.0-20191209.1573068157933.x86_64
func (p *Pkg) GetNVRA() string {
	release := "None"
	arch := "None"
	if len(p.Release) > 0 {
		release = p.Release
	}
	if len(strings.Split(p.Platform, "-")) > 0 {
		arch = strings.Split(p.Platform, "-")[0]
	}
	return p.Name + "-" + p.Version + "-" + release + "." + arch
}

// FetchSources fetches the sources from a pkg
func (p *Pkg) FetchSources(destdir string, force bool) ([]string, error) {
	filelist := make([]string, 0)
	err := os.MkdirAll(destdir, 0755)
	if err != nil {
		return nil, err
	}
	for _, src := range p.Sources {
		fmt.Printf("Archive : %s\n", src.Archive)
		filename := path.Base(src.Archive)
		filepath := path.Join(destdir, filename)
		fmt.Printf("FilePath : %s\n", filepath)
		if force || !common.IsFile(filepath) {
			if common.IsDir(filepath) {
				return nil, err
			}
			err := files.DownloadFile(filepath, src.Archive)
			if err != nil {
				return nil, err
			}
		}

		raw, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		md5sum := fmt.Sprintf("%x", md5.Sum(raw)) // nolint:gosec
		if src.MD5 != "" {
			if md5sum != src.MD5 {
				return nil, fmt.Errorf("%s : MD5 sums do not match %s != %s", src.Archive, src.MD5, md5sum)
			}
		}
		sha256sum := fmt.Sprintf("%x", sha256.Sum256(raw))
		if src.SHA256 != "" {
			if sha256sum != src.SHA256 {
				return nil, fmt.Errorf("%s : SHA256 sums do not match %s != %s", src.Archive, src.SHA256, sha256sum)
			}
		}
		filelist = append(filelist, sha256sum+" "+filepath)
		f := files.File{
			Path:   filepath,
			Name:   filename,
			Mode:   "0644", // TODO: get actual mode
			SHA256: sha256sum,
			MD5:    md5sum,
		}
		p.Files = append(p.Files, f)
	}
	return filelist, nil
}

// ToYAML func takes no input and returns []byte, error
func (p *Pkg) ToYAML() ([]byte, error) {
	content, err := yaml.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON func takes no input and returns []byte, error
func (p *Pkg) ToJSON() ([]byte, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}

// ToPrettyJSON func takes no input and returns []byte, error
func (p *Pkg) ToPrettyJSON() ([]byte, error) {
	content, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}

// ToBuildScript func takes no input and returns []byte, error
func (p *Pkg) ToBuildScript() (string, error) {
	template := `#!/bin/bash
set -x
umask 022
LANG=C
LC_ALL=POSIX
PATH=/tools/bin:/bin:/usr/bin
export LANG LC_ALL PATH
PKG_CONFIG_PATH="${PKG_CONFIG_PATH}:/usr/lib64/pkgconfig:/usr/share/pkgconfig"
export PKG_CONFIG_PATH
BUILDDIR=/build
SRCDIR=/src
DESTDIR=/install
PKGDIR=/package
export BUILDDIR SRCDIR DESTDIR PKGDIR
mkdir -vp {$BUILDDIR,$SRCDIR,$DESTDIR,$PKGDIR}
{{ range .Instructions}}
{{.Unpack}}
{{.Pre}}
{{.Configure}}
{{.Build}}
{{.Test}}
{{.Install}}
{{.Post}}
{{end}}
`
	data := structs.Map(p)
	script, err := common.MakeTemplate(data, template)
	if err != nil {
		return script, err
	}
	return script, nil
}

// ToYAML func takes no input and returns []byte, error
func (p *Pkgs) ToYAML() ([]byte, error) {
	content, err := yaml.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON func takes no input and returns []byte, error
func (p *Pkgs) ToJSON() ([]byte, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}

// DecodePkgFromJSON func takes reader io.Reader as input and returns *Pkg, error
func DecodePkgFromJSON(reader io.Reader) (*Pkg, error) {
	pkg := &Pkg{}
	err := json.NewDecoder(reader).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("error decoding json : %v", err)
	}
	return pkg, nil
}

// DecodePkgFromYAML func takes reader io.Reader as input and returns *Pkg, error
func DecodePkgFromYAML(reader io.Reader) (*Pkg, error) {
	pkg := &Pkg{}
	err := yaml.NewDecoder(reader).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("error decoding yaml : %v", err)
	}
	return pkg, nil
}

// Build attempts to build the package
// Make Build Directory
// Create Build Script
// Run Build script
// Compress Results
func (p *Pkg) Build(builddir string) (string, error) {
	var output string
	err := os.MkdirAll(builddir, 0755)
	if err != nil {
		return output, err
	}

	return output, nil
}
