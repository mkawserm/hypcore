package app

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

func CheckConfigFile(configFilePath string) *viper.Viper {
	v := NewConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if configFilePath != "" {
				fmt.Println(aurora.BrightRed("config file not found @ " + configFilePath))
			} else {
				fmt.Println(aurora.BrightRed("config file not found @ " + ConfigFilePathFirst + "/" + ConfigFileName()))
			}
		} else {
			errorStr := fmt.Sprintf("config file parsing error: %s", err.Error())
			fmt.Printf("%s", aurora.BrightRed(errorStr))
		}
	}

	return v
}
