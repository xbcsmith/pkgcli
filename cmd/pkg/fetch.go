// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes" // nolint:gosec
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbcsmith/pkgcli/models"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch package binaries",
	Long:  `fetch package binaries`,
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
		sourcedir := viper.GetString("sourcedir")
		force := viper.GetBool("force")
		fmt.Printf("SOURCEDIR : %s\n", sourcedir)
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
				fmt.Printf("Name : %s\n", pkg.Name)
				fmt.Printf("Version : %s\n", pkg.Version)
				fmt.Printf("Release : %s\n", pkg.Release)
				if err := os.MkdirAll(sourcedir, 0755); err != nil {
					return err
				}
				var filelist []string

				fmt.Printf("Downloading to %s\n", sourcedir)
				filelist, err = pkg.FetchSources(sourcedir, force)
				if err != nil {
					return err
				}

				fmt.Println("Fetch complete")
				for _, file := range filelist {
					fmt.Printf("%s\n", file)
				}
			}

		}
		return nil
	},
}

// NewFetchCmd returns a new fetch command
func NewFetchCmd() *cobra.Command {
	fetchCmd.Flags().String("sourcedir", "/src", "Source Directory directory")
	fetchCmd.Flags().Bool("force", false, "Force Fetch Sources even if in sourcedir")

	return fetchCmd
}
