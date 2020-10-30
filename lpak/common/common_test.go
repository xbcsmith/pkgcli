// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

// TestStringInSlice validates StringInSlice works
// TestStringInSlice func takes no input and returns t *testing.T
func TestStringInSlice(t *testing.T) {
	list := []string{"foo", "bar", "caz"}
	result := StringInSlice("foo", list)
	assert.Equal(t, result, true, "Error in StringInSlice")
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
