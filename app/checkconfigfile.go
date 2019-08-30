package app

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

func CheckConfigFile(configFilePath string, configFileName string) *viper.Viper {
	v := NewConfigFile(configFilePath, configFileName)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if configFilePath != "" {
				fmt.Println(aurora.BrightRed("config file not found @ " + configFilePath))
			} else {
				fmt.Println(aurora.BrightRed("config file not found @ " + ConfigFilePathFirst + "/" + ConfigFileName()))
			}
			return nil
		} else {
			errorStr := fmt.Sprintf("config file parsing error: %s", err.Error())
			fmt.Printf("%s", aurora.BrightRed(errorStr))
			return nil
		}
	}

	return v
}

// Check if the configuration file is ok
func IsConfigurationOk(v *viper.Viper) bool {

	return false
}
