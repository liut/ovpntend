package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ovpntend",
	Short: "ovpntend and command launcher " + version,
	Long:  ``,
}

func Execute() {
	if inDevelop() {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if inDevelop() {
		slog.Info("in develop")
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	startUp()
}
