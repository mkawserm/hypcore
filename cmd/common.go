package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/mkawserm/hypcore/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
	"strings"
)

func Common(cmd *cobra.Command, args []string, logToGlog bool) *viper.Viper {
	filePath := ""
	fileNameWithoutExt := ""

	file, err := cmd.Flags().GetString("config")

	if err == nil {
		if file != "" {
			filePath, fileNameWithoutExt = path.Split(file)
			if strings.ToLower(path.Ext(fileNameWithoutExt)) != ".toml" {
				if logToGlog {
					glog.Errorln("Configuration file with .toml extension is allowed")
					glog.Errorln(path.Ext(fileNameWithoutExt) + " extension is not allowed")
				} else {
					fmt.Println("Configuration file with .toml extension is allowed")
					fmt.Println(path.Ext(fileNameWithoutExt) + " extension is not allowed")
				}
				return nil
			}
			fileNameWithoutExt = strings.TrimSuffix(fileNameWithoutExt, path.Ext(fileNameWithoutExt))
		} else {
			filePath = app.ConfigFilePathFirst
			fileNameWithoutExt = app.ConfigFileNameWithoutExt
		}

		v := app.CheckConfigFile(filePath, fileNameWithoutExt, false, false)
		if v == nil && !app.IsConfigurationOk(v, false, false) {
			if logToGlog {
				glog.Errorln("Invalid configuration file")
			} else {
				fmt.Println("Invalid configuration file")
			}

			return nil
		}
		return v
	}

	return nil
}
