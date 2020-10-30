// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbcsmith/pkgcli/cmd/install"
	"github.com/xbcsmith/pkgcli/cmd/pkg"
	"github.com/xbcsmith/pkgcli/cmd/remove"
)

// rootCmd for cobra
var rootCmd = &cobra.Command{
	Use:   "pkgcli",
	Short: "Command Line Package helper",
	Long:  `Command Line utility for managing package installs`,
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
	RunE: run,
}

// Execute runs things
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println(GetVersion())
	fmt.Printf("\nusage: " + AppName + " --help\n\n")
	return nil
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debugging statements")

	// install
	installCmd := install.NewInstallCmd()
	// build
	removeCmd := remove.NewRemoveCmd()
	// pkgs
	pkgCmd := pkg.NewPkgCmd()
	pkgCreate := pkg.NewCreateCmd()
	pkgBuild := pkg.NewBuildCmd()
	pkgFetch := pkg.NewFetchCmd()

	pkgCmd.AddCommand(pkgCreate)
	pkgCmd.AddCommand(pkgBuild)
	pkgCmd.AddCommand(pkgFetch)

	// Add commands to root
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(pkgCmd)
}
