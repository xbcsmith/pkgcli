// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd for cobra
var RootCmd = &cobra.Command{
	Use:   "lfscli",
	Short: "Command Line LFS helper",
	Long:  `Command Line LFS utility for managing package installs`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
	if err := RootCmd.Execute(); err != nil {
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
	RootCmd.PersistentFlags().Bool("debug", false, "Enable debugging statements")

	// pkgs
	pkgCmd := pkg.NewPkgCmd()
	pkgCreate := pkg.NewCreateCmd()
	pkgBuild := pkg.NewQueryCmd()
	pkgSearch := pkg.NewSearchCmd()
	pkgCmd.AddCommand(pkgCreate)
	pkgCmd.AddCommand(pkgBuild)
	pkgCmd.AddCommand(pkgSearch)

	//Add commands to root
	RootCmd.AddCommand(pkgCmd)

}
