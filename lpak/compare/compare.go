// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package compare

import (
	"github.com/oklog/ulid"
	"github.com/xbcsmith/pkgcli/lpak/model"
)

// TODO: write a version compare... man why?
// You know what last in last out. ULID are sortable.
// I will only sort by ULID and return the latest release.
// by not considering version we eliminate the epoch situation
// release is the latest
const (
	Lesser  = -1
	Equal   = 0
	Greater = 1
)

func Compare(x, y *model.Pkg) int {
	a := ulid.MustParse(x.Release)
	b := ulid.MustParse(y.Release)
	return a.Compare(b)
}
