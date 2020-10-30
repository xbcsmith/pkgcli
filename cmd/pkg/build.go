// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbcsmith/pkgcli/lpak/common"
	"github.com/xbcsmith/pkgcli/lpak/model"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build package binaries",
	Long:  `build package binaries`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("PACKAGE")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		buildroot := viper.GetString("buildroot")
		sourcedir := viper.GetString("sourcedir")
		force := viper.GetBool("force-download")
		fmt.Printf("BUILDROOT : %s\nSOURCEDIR : %s\n", buildroot, sourcedir)
		if force {
			fmt.Println("Force Enabled")
		}
		if len(args) > 0 {
			for _, filepath := range args {
				content, err := ioutil.ReadFile(filepath)
				if err != nil {
					return err
				}
				var p = &model.Pkg{}
				if !common.IsJSON(content) {
					p, err = model.DecodePkgFromYAML(bytes.NewReader(content))
					if err != nil {
						return err
					}
				} else {
					p, err = model.DecodePkgFromJSON(bytes.NewReader(content))
					if err != nil {
						return err
					}
				}
				if len(p.Release) == 0 {
					p.Release = common.NewRelease()
				}
				fmt.Printf("NVRA : %s\n", p.GetNVRA())

				artifacts, err := p.FetchSources(sourcedir, force)
				if err != nil {
					return nil
				}
				fmt.Println(common.SHASlice(artifacts))
			}

		}
		return nil
	},
}

// NewBuildCmd returns a new build command
func NewBuildCmd() *cobra.Command {
	buildCmd.Flags().String("buildroot", "/tmp", "Build Root directory")
	buildCmd.Flags().String("sourcedir", "/src", "Source Directory directory")
	buildCmd.Flags().Bool("fetch", false, "Fetch Sources if not in sourcedir")

	return buildCmd
}
