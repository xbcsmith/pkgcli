// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/xbcsmith/pkgcli/lpak/model"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestCreateDB(t *testing.T) {
	err := CreateDB("/tmp/db.sql")
	assert.Assert(t, is.Nil(err))
}

func TestConnect(t *testing.T) {
	db := Connect("/tmp/db.sql")
	columnNames(db)
	cleanUp("/tmp/db.sql")
}

func cleanUp(path string) {
	e := os.Remove(path)
	if e != nil {
		log.Fatal(e)
	}
}

func columnNames(db *gorm.DB) {
	fmt.Println("============================")
	fmt.Println("Column Names")
	fmt.Println("============================")

	for _, f := range db.NewScope(&model.Pkg{}).Fields() {
		println(f.Name)
	}
}
