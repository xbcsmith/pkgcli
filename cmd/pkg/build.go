// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbcsmith/pkgcli/models"
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
					fmt.Println(err)
					os.Exit(-1)
				}
				pkg := &models.Pkg{
					Description:  "",
					Instructions: []models.Instruction{},
					Name:         "",
					Package:      "",
					Platform:     "",
					Provides:     []string{},
					Release:      "",
					Requires:     []string{},
					Optional:     []string{},
					Recommended:  []string{},
					Sources:      []models.Source{},
					Files:        []models.File{},
					Summary:      "",
					Version:      "",
				}
				isjson := models.IsJSON(content)
				if !isjson {
					pkg, err = models.DecodePkgFromYAML(bytes.NewReader(content))
					if err != nil {
						return err
					}
				} else {
					pkg, err = models.DecodePkgFromJSON(bytes.NewReader(content))
					if err != nil {
						return err
					}
				}
				if len(pkg.Release) == 0 {
					pkg.Release = models.NewRelease()
				}
				fmt.Printf("NVRA : %s\n", pkg.GetNVRA())

				artifacts, err := pkg.FetchSources(sourcedir, force)
				if err != nil {
					return nil
				}
				fmt.Println(models.SHASlice(artifacts))
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
