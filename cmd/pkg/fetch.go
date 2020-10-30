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
	"github.com/xbcsmith/pkgcli/lpak/common"
	"github.com/xbcsmith/pkgcli/lpak/model"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch package binaries",
	Long:  `fetch package binaries`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("package")
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
				fmt.Printf("Name : %s\n", p.Name)
				fmt.Printf("Version : %s\n", p.Version)
				fmt.Printf("Release : %s\n", p.Release)
				if err := os.MkdirAll(sourcedir, 0755); err != nil {
					return err
				}
				var filelist []string

				fmt.Printf("Downloading to %s\n", sourcedir)
				filelist, err = p.FetchSources(sourcedir, force)
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
