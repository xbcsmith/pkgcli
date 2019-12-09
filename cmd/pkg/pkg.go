package pkg

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pkgCmd represents the pkg command
var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "pkg command",
	Long:  `pkg command`,
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

		return nil
	},
}

// NewPkgCmd returns a new pkg command
func NewPkgCmd() *cobra.Command {
	return pkgCmd
}
