package pkg

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbcsmith/pkgcli/models"
	"github.com/xbcsmith/pkgcli/utils"
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
				pkg := &models.Pkg{}
				isjson := utils.IsJSON(content)
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
				var filelist []string
				if !force {
					for _, src := range pkg.Sources {
						filelist, err = utils.FindFile(sourcedir, path.Base(src.Archive))
						for _, filepath := range filelist {
							raw, err := ioutil.ReadFile(filepath)
							if err != nil {
								fmt.Println(err)
								os.Exit(-1)
							}
							md5sum := fmt.Sprintf("%x", md5.Sum(raw))
							if md5sum != src.MD5 {
								fmt.Printf("%s : MD5 sums do not match %s != %s", src.Archive, src.MD5, md5sum)
							}
							sha256sum := fmt.Sprintf("%x", sha256.Sum256(raw))
							if sha256sum != src.SHA256 {
								fmt.Printf("%s : SHA256 sums do not match %s != %s", src.Archive, src.SHA256, sha256sum)
							}
						}
					}
				} else {
					filelist, err = pkg.FetchSources(sourcedir)
					if err != nil {
						return err
					}
				}
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
