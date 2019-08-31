package cmd

import (
	"github.com/golang/glog"
	"github.com/mkawserm/hypcore/app"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

func RunServerCmdRun(cmd *cobra.Command, args []string) {
	filePath := ""
	fileNameWithoutExt := ""

	file, err := cmd.Flags().GetString("config")

	if err == nil {
		if file != "" {
			filePath, fileNameWithoutExt = path.Split(file)
			if strings.ToLower(path.Ext(fileNameWithoutExt)) != ".toml" {
				glog.Errorln("Configuration file with .toml extension is allowed")
				glog.Errorln(path.Ext(fileNameWithoutExt) + " extension is not allowed")
				return
			}
			fileNameWithoutExt = strings.TrimSuffix(fileNameWithoutExt, path.Ext(fileNameWithoutExt))
		} else {
			filePath = app.ConfigFilePathFirst
			fileNameWithoutExt = app.ConfigFileNameWithoutExt
		}

		v := app.CheckConfigFile(filePath, fileNameWithoutExt, false, false)
		if v == nil && !app.IsConfigurationOk(v, false, false) {
			glog.Errorln("Invalid configuration file")
			return
		}

		app.RunServer(v)
	}
}
