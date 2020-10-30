package compare

import (
	"fmt"
	"testing"

	"github.com/xbcsmith/pkgcli/lpak/files"
	"github.com/xbcsmith/pkgcli/lpak/instructions"
	"github.com/xbcsmith/pkgcli/lpak/model"
	"github.com/xbcsmith/pkgcli/lpak/source"
)

func TestCompare(t *testing.T) {
	a := &model.Pkg{
		Description:  "",
		Instructions: []instructions.Instruction{},
		Name:         "",
		Package:      "",
		Platform:     "",
		Provides:     []string{},
		Release:      "",
		Requires:     []string{},
		Optional:     []string{},
		Recommended:  []string{},
		Sources:      []source.Source{},
		Files:        []files.File{},
		Summary:      "",
		Version:      "",
	}
	fmt.Printf("%+v\n", a)
}
