package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"fhyx.tech/platform/ovpntend/pkg/zlog"
)

var RootCmd = &cobra.Command{
	Use:   "ovpntend",
	Short: "ovpntend and command launcher " + version,
	Long:  ``,
}

func Execute() {
	var zlogger *zap.Logger
	if inDevelop() {
		zlogger, _ = zap.NewDevelopment()
	} else {
		zlogger, _ = zap.NewProduction()
	}
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()

	zlog.Set(sugar)
	if inDevelop() {
		logger().Infow("in develop")
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	startUp()
}

func logger() zlog.Logger {
	return zlog.Get()
}
