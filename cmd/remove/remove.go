// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package remove

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove command",
	Long:  `remove command`,
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
		fmt.Println("This is the remove command")
		return nil
	},
}

// NewRemoveCmd returns a new remove command
func NewRemoveCmd() *cobra.Command {
	return removeCmd
}
