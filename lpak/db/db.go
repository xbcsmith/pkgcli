// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/jinzhu/gorm"
	"github.com/xbcsmith/pkgcli/lpak/model"
)

// Config is the configuration for the db
type Config struct {
	Path string
}

// CreateDB creates db
func CreateDB(path string) error {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&model.Pkg{})
	return nil
}

// Connect connects to the db
func Connect(path string) *gorm.DB {
	// Test if path exists
	// Test if path is writable
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	return db
}

// Create a record in the database
func Create(pkg *model.Pkg) error {
	return nil
}

// Update a record in the database
func Update(pkg *model.Pkg) error {
	return nil
}

// Delete a record in the database
func Delete(pkg *model.Pkg) error {
	return nil
}
