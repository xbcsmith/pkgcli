// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"testing"

	"gotest.tools/v3/assert"
)

// TestStringInSlice validates StringInSlice works
// TestStringInSlice func takes no input and returns t *testing.T
func TestStringInSlice(t *testing.T) {
	list := []string{"foo", "bar", "caz"}
	result := StringInSlice("foo", list)
	assert.Equal(t, result, true, "Error in StringInSlice")
}
