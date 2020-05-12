package cmd

import (
	"github.com/spf13/cobra"

	"fhyx.tech/platform/ovpntend/pkg/settings"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Print usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		settings.Usage()
	},
}

func init() {
	RootCmd.AddCommand(usageCmd)

}
