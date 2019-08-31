package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/app"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
	"strings"
)

func CreateConfigFileCmdRun(cmd *cobra.Command, args []string) {
	filePath := ""
	fileNameWithoutExt := ""

	var e error
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
		e = app.EnsureDir(filePath)
		if e != nil {
			fmt.Println(aurora.BrightRed(e.Error()))
			return
		}

		v := app.NewConfigFile(filePath, fileNameWithoutExt)
		_ = ioutil.WriteFile(filePath+"/"+fileNameWithoutExt+".toml", []byte(""), 0777)
		e = v.WriteConfig()
		if e != nil {
			fmt.Println(aurora.BrightRed(e.Error()))
			return
		} else {
			fmt.Println(aurora.BrightGreen("Configuration file created"))
			return
		}
	}
}
