package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	built   = "N/A"
	name    = "ovpntend"
	version = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func inDevelop() bool {
	return version == "dev"
}

func printVersion() {
	fmt.Printf("%s %s built %s (%s %s-%s)\n", name, version, built, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
