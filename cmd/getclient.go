// Copyright © 2020 liut <liutao@liut.cc>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"io/ioutil"
	"log/slog"
	"path"

	"github.com/spf13/cobra"

	"github.com/liut/ovpntend/pkg/ovpn"
)

// getclientCmd represents the getclient command
var getclientCmd = &cobra.Command{
	Use:   "getclient",
	Short: "Get a client config and save it into a directory",
	Run:   getclientRun,
}

func init() {
	RootCmd.AddCommand(getclientCmd)

	getclientCmd.Flags().StringP("name", "n", "", "A special client name")
	getclientCmd.Flags().StringP("out", "w", "", "A directory for output")
	getclientCmd.Flags().BoolP("sendmail", "s", false, "send a email to client name")
	getclientCmd.Flags().String("os", "mac", "OS category")
}

func getclientRun(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	out, _ := cmd.Flags().GetString("out")
	sendmail, _ := cmd.Flags().GetBool("sendmail")
	oscat, _ := cmd.Flags().GetString("os")

	if 0 == len(name) || 0 == len(out) && !sendmail {
		slog.Info("empty name or output directory")
		cmd.Usage()
		return
	}

	ctx := context.Background()
	if sendmail {
		if err := ovpn.SendConfig(ctx, name, oscat); err != nil {
			slog.Warn("sendmail fail", "err", err)
		} else {
			slog.Info("sendmail OK")
		}
		return
	}

	body, err := ovpn.GetClientConfig(ctx, name)
	if err != nil {
		slog.Info("get fail", "err", err)
		return
	}

	file := path.Join(out, name+".ovpn")
	err = ioutil.WriteFile(file, body, 0644)
	if err != nil {
		slog.Info("write fail", "err", err)
		return
	}
	slog.Info("saved ok", "file", file)
}
