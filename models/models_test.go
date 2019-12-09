package models

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestDecodePkgFromJson(t *testing.T) {
	content := `{
  "description": "Mech from Brixton",
  "id": "01DM3PADN0FQFGD2SDJR04DGF8",
  "name": "foo",
  "release": "20190906.1567787180",
  "version": "0.0.1"
} `

	mech, err := DecodePkgFromJSON(strings.NewReader(content))
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, mech.Name, "foo")

}

func TestNewPkg(t *testing.T) {
	mech := NewPkg()
	assert.Assert(t, len(mech.ID) == 26)
}
