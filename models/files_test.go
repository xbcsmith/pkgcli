// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"path"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

const sharutilFileName string = "sharutils.yaml"

func TestFindDir(t *testing.T) {
	dirname := "pkg"
	basepath := "../testdata/"
	expected := "pkg"
	dirs, err := FindDir(basepath, dirname, 2)
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, path.Base(dirs[0]), expected, "Expected text failed")
}

func TestFindFile(t *testing.T) {
	filename := sharutilFileName
	filepath := "../testdata/pkg/"
	expected := sharutilFileName
	files, err := FindFile(filepath, filename)
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, path.Base(files[0]), expected, "Expected text failed")
}

func TestFindArtifacts(t *testing.T) {
	dir := "../testdata/pkg"
	suffix := ".yaml"
	expected := sharutilFileName
	artifacts, err := FindArtifacts(dir, suffix)
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, path.Base(artifacts[0]), expected, "Expected paths failed")
}
