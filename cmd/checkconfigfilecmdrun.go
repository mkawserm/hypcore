package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/app"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

func CheckConfigFileCmdRun(cmd *cobra.Command, args []string) {
	filePath := ""
	fileNameWithoutExt := ""

	file, err := cmd.Flags().GetString("file")
	if err == nil {
		if file != "" {
			filePath, fileNameWithoutExt = path.Split(file)
			if strings.ToLower(path.Ext(fileNameWithoutExt)) != ".toml" {
				fmt.Println(aurora.BrightRed("Configuration file with .toml extension is allowed"))
				fmt.Println(aurora.BrightRed(path.Ext(fileNameWithoutExt) + " extension is not allowed"))
				return
			}

			fileNameWithoutExt = strings.TrimSuffix(fileNameWithoutExt, path.Ext(fileNameWithoutExt))
		} else {
			filePath = app.ConfigFilePathFirst
			fileNameWithoutExt = app.ConfigFileNameWithoutExt
		}

		v := app.CheckConfigFile(filePath, fileNameWithoutExt)
		if v != nil && app.IsConfigurationOk(v, false) {
			fmt.Println(aurora.BrightGreen("Configuration file Ok"))
			return
		} else {
			fmt.Println(aurora.BrightRed("Invalid configuration file"))
			return
		}
	}
}
