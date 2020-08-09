// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"

	"crypto/md5" // nolint:gosec
	"crypto/sha256"

	"github.com/spf13/cobra"
	"github.com/xbcsmith/pkgcli/models"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an lfs pkg",
	Long:  `Creates an lfs pkg yaml`,
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

		name := viper.GetString("name")
		version := viper.GetString("version")
		release := viper.GetString("release")
		description := viper.GetString("description")
		summary := viper.GetString("summary")
		pkgType := viper.GetString("package")
		platformID := viper.GetString("platform_id")
		provides := viper.GetStringSlice("provides")
		requires := viper.GetStringSlice("requires")

		pkg := models.NewPkg(name, version)

		//TODO: validate parameters
		if release != "" {
			pkg.Release = release
		}
		if description != "" {
			pkg.Description = description
		}
		if summary != "" {
			pkg.Summary = summary
		}
		if pkgType != "" {
			pkg.Package = pkgType
		}
		if platformID != "" {
			pkg.PlatformID = platformID
		}
		pkg.Provides = provides
		pkg.Requires = requires

		if len(args) > 0 {
			for _, filepath := range args {
				raw, err := ioutil.ReadFile(filepath)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				md5sum := md5.Sum(raw) // nolint:gosec
				sha256sum := sha256.Sum256(raw)
				src := &models.Source{
					Archive: filepath,
					MD5:     fmt.Sprintf("%x", md5sum),
					SHA256:  fmt.Sprintf("%x", sha256sum),
				}
				pkg.Sources = append(pkg.Sources, *src)
				unpack := "tar -xvf " + filepath + " && cd " + name + "-" + version
				instruction := &models.Instruction{
					Build:     "make",
					Configure: "./configure --prefix=/usr",
					Install:   "make install",
					Post:      "",
					Pre:       "",
					Test:      "make check",
					Unpack:    unpack,
				}
				pkg.Instructions = append(pkg.Instructions, *instruction)
			}
		}

		yml, err := pkg.ToYAML()
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", yml)

		return nil
	},
}

// NewCreateCmd creates a new cmdline
func NewCreateCmd() *cobra.Command {
	createCmd.Flags().String("name", "", "Name of the pkg to create")
	createCmd.Flags().String("release", "", "Release of the pkg to create")
	createCmd.Flags().String("version", "", "Version of the pkg to create")
	createCmd.Flags().String("description", "", "Description of the pkg to create")
	createCmd.Flags().String("summary", "", "Summary of the pkg to create")

	createCmd.Flags().String("package", "", "Package Type of the pkg to create (tar.xz, tar.gz, tgz)")
	createCmd.Flags().String("platform_id", "", "Platform ID of the pkg to create (x86_64-gnu-linux-9)")
	createCmd.Flags().String("requires", "", "Requires of the pkg to create (bar, caz)")
	createCmd.Flags().String("provides", "", "Provides of the pkg to create (libfoo.so.1)")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("version")

	return createCmd
}
