package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wasilak/notes-manager/libs/common"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + common.AppName,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.SetContext(common.CTX)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := versionFunc(); err != nil {
			return err
		}
		return nil
	},
}

func versionFunc() error {

	fmt.Printf("%s\nVersion %s (GO %s)\n", common.AppName, common.Version, common.GetVersion())
	return nil
}
