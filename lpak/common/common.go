// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"

	//nolint:gosec
	"crypto/rand"

	"github.com/oklog/ulid"
)

// SHASlice reps a sha256sum file
type SHASlice []string

// String returns a sha256sum file format
func (s SHASlice) String() string {
	var str string
	for _, s := range s {
		str += fmt.Sprintf("%s\n", s)
	}
	return str
}

// NewULID returns a ULID.
func NewULID() (ulid.ULID, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return id, fmt.Errorf("NewULID Failed: %s", err)
	}
	return id, err
}

// NewULIDAsString returns a ULID string.
func NewULIDAsString() string {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}

// NewRelease returns a ULID as a release.
// using ULID as a release solves a crazy compare
// problem and removes the need for an epoch
func NewRelease() string {
	return NewULIDAsString()
}

// MakeTemplate helper function
func MakeTemplate(data map[string]interface{}, tmpl string) (string, error) {
	builder := &strings.Builder{}
	t := template.Must(template.New("new").Parse(tmpl))
	if err := t.Execute(builder, data); err != nil {
		return ``, err
	}
	s := builder.String()
	return s, nil
}

// StringInSlice is used to find a string in a list
// StringInSlice func takes a string, array and returns a bool
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IsJSON try to guess if file is JSON or YAML.
func IsJSON(buf []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	if bytes.HasPrefix(trim, []byte("{")) {
		return true
	}
	if bytes.HasPrefix(trim, []byte("[")) {
		return true
	}
	return false
}

// IsFile returns true if file exists and is not a dir
func IsFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDir returns true if dir exists
func IsDir(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
