package pkg

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx/types"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an lfs pkg",
	Long:  `Creates an lfs pkg yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
    package := viper.GetString("package")
    platformID := viper.GetString("platform_id")
    provides := viper.GetString("provides")
    requires := viper.GetString("requires")


		//TODO: validate parameters

		pkg := &pkg.Pkg{
    	Name: name,
      Version: version,
      Release: release,
      Description: description,
      Summary: summary,
    	Package: package,
    	PlatformID: platformID,
    	Provides: provides,
    	Requires: requires,
    	Sources: sources,
    	Instructions: instructions,
		}

    fmt.Printf("%s\n", utils.ToYaml(pkg))

		return nil
	},
}

// NewCreateCmd creates a new cmdline
func NewCreateCmd() *cobra.Command {
	createCmd.Flags().String("name", "", "A name for the pkg")
	createCmd.Flags().String("action", "", "An action for the pkg")
	createCmd.Flags().String("version", "", "Version of the pkg to create")
	createCmd.Flags().String("description", "", "Description of the pkg")
	createCmd.Flags().String("schema-version", "", "What version of the schema you are attaching")
	createCmd.Flags().String("maintainer", "", "Name and email of maintainer")
	createCmd.Flags().String("tags", "", "Space delimited list of tags")
	createCmd.Flags().String("encoding", "utf-8", "Encoding you're using")
	createCmd.Flags().String("schema", "{}", "Schema to validate receipts passing through")
	createCmd.Flags().String("ca-cert", "", "TLS Certificate")
	createCmd.Flags().String("metadata", "{}", "Metadata about the pkg")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("action")
	_ = createCmd.MarkFlagRequired("version")
	_ = createCmd.MarkFlagRequired("description")
	_ = createCmd.MarkFlagRequired("schema-version")
	_ = createCmd.MarkFlagRequired("maintainer")
	_ = createCmd.MarkFlagRequired("tags")
	_ = createCmd.MarkFlagRequired("schema")
	_ = createCmd.MarkFlagRequired("metadata")

	return createCmd
}
