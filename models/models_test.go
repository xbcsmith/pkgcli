// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

// TestDecodePkgFromYAML func takes no input and returns t *testing.T
func TestDecodePkgFromYAML(t *testing.T) {
	content := `name: sharutils
version: 4.15.2
release: null
description: "The Sharutils package contains utilities that can create 'shell' archives"
summary: "sharutils utilities for shell archives"
requires: []
provides: []
instructions:
- unpack: tar -xvf $name-$version.tar.xz && cd $name-$version
  pre: "sed -i 's/IO_ftrylockfile/IO_EOF_SEEN/' lib/*.c && echo '#define _IO_IN_BACKUP 0x100' >> lib/stdio-impl.h"
  configure: "./configure --prefix=/usr"
  build: "make"
  test: "make check"
  install: "make install"
  post: ""
sources:
- archives: http://ftp.gnu.org/gnu/sharutils/sharutils-4.15.2.tar.xz
  md5: 5975ce21da36491d7aa6dc2b0d9788e0
  sha256: 2b05cff7de5d7b646dc1669bc36c35fdac02ac6ae4b6c19cb3340d87ec553a9a
  destination: /usr
platform_id: x86_64-lfs-linux-9
package: tar.xz
`

	pkg, err := DecodePkgFromYAML(strings.NewReader(content))
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, pkg.Name, "sharutils")
	arch := strings.Split(pkg.PlatformID, "-")[0]
	assert.Equal(t, arch, "x86_64")
}

// TestDecodePkgFromJSON func takes no input and returns t *testing.T
func TestDecodePkgFromJSON(t *testing.T) {
	content := `{
  "description": "The Sharutils package contains utilities that can create 'shell' archives",
  "instructions": [{
    "build": "make",
    "configure": "./configure --prefix=/usr",
    "install": "make install",
    "post": "",
    "pre": "sed -i 's/IO_ftrylockfile/IO_EOF_SEEN/' lib/*.c \u0026\u0026 echo '#define _IO_IN_BACKUP 0x100' \u003e\u003e lib/stdio-impl.h",
    "test": "make check",
    "unpack": "tar -xvf $name-$version.tar.xz \u0026\u0026 cd $name-$version"
  }],
  "name": "sharutils",
  "package": "tar.xz",
  "platform_id": "x86_64-lfs-linux-9",
  "provides": [],
  "release": null,
  "requires": [],
  "sources": [{
    "archives": "http://ftp.gnu.org/gnu/sharutils/sharutils-4.15.2.tar.xz",
    "destination": "/usr",
    "md5": "5975ce21da36491d7aa6dc2b0d9788e0",
    "sha256": "2b05cff7de5d7b646dc1669bc36c35fdac02ac6ae4b6c19cb3340d87ec553a9a"
  }],
  "summary": "sharutils utilities for shell archives",
  "version": "4.15.2"
}`

	pkg, err := DecodePkgFromJSON(strings.NewReader(content))
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, pkg.Name, "sharutils")
}

// TestNewPkg func takes no input and returns t *testing.T
func TestNewPkg(t *testing.T) {
	pkg := NewPkg("foo", "2.1.1")
	assert.Assert(t, pkg.Name == "foo")
	assert.Assert(t, pkg.Version == "2.1.1")
}

// TestNewULID func takes no input and returns t *testing.T
func TestNewULID(t *testing.T) {
	u, err := NewULID()
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, len(u.String()) == 26)
}

// TestNewULIDAsString func takes no input and returns t *testing.T
func TestNewULIDAsString(t *testing.T) {
	u := NewULIDAsString()
	assert.Assert(t, len(u) == 26)
}

// TestDecodePkgFromYAML func takes no input and returns t *testing.T
func TestPkgToBuildScript(t *testing.T) {
	content := `name: sharutils
version: 4.15.2
release: null
description: "The Sharutils package contains utilities that can create 'shell' archives"
summary: "sharutils utilities for shell archives"
requires: []
provides: []
instructions:
- unpack: tar -xvf $SRCDIR/sharutils-4.15.2.tar.xz && cd sharutils-4.15.2
  pre: "sed -i 's/IO_ftrylockfile/IO_EOF_SEEN/' lib/*.c && echo '#define _IO_IN_BACKUP 0x100' >> lib/stdio-impl.h"
  configure: "./configure --prefix=/usr"
  build: "make"
  test: "make check"
  install: "make install DESTDIR=$DESTDIR"
  post: ""
sources:
- archives: http://ftp.gnu.org/gnu/sharutils/sharutils-4.15.2.tar.xz
  md5: 5975ce21da36491d7aa6dc2b0d9788e0
  sha256: 2b05cff7de5d7b646dc1669bc36c35fdac02ac6ae4b6c19cb3340d87ec553a9a
  destination: /usr
platform_id: x86_64-lfs-linux-9
package: tar.xz
`

	pkg, err := DecodePkgFromYAML(strings.NewReader(content))
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, pkg.Name, "sharutils")
	script, err := pkg.ToBuildScript()
	assert.Assert(t, is.Nil(err))
	fmt.Printf("%s\n", script)
}
