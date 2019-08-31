package app

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

func CheckConfigFileSilent(configFilePath string, configFileName string) *viper.Viper {
	v := NewConfigFile(configFilePath, configFileName)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//if configFilePath != "" {
			//	fmt.Println(aurora.BrightRed("config file not found @ " + configFilePath))
			//} else {
			//	fmt.Println(aurora.BrightRed("config file not found @ " + ConfigFilePathFirst + "/" + ConfigFileName()))
			//}
			return v
		} else {
			//errorStr := fmt.Sprintf("config file parsing error: %s", err.Error())
			//fmt.Printf("%s", aurora.BrightRed(errorStr))
			return nil
		}
	}

	return v
}

func CheckConfigFile(configFilePath string,
	configFileName string,
	silent bool, logToGlog bool) *viper.Viper {
	v := NewConfigFile(configFilePath, configFileName)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if configFilePath != "" {
				if silent == false {
					if logToGlog == true {
						glog.Errorln("config file not found @ " + configFilePath + "/" + configFileName + ".toml")
					} else {
						fmt.Println(aurora.BrightRed("config file not found @ " + configFilePath + "/" + configFileName + ".toml"))
					}
				}

			} else {
				if silent == false {
					if logToGlog == true {
						glog.Errorln("config file not found @ " + ConfigFilePathFirst + "/" + ConfigFileName())
					} else {
						fmt.Println(aurora.BrightRed("config file not found @ " + ConfigFilePathFirst + "/" + ConfigFileName()))
					}
				}
			}
			return nil
		} else {
			errorStr := fmt.Sprintf("config file parsing error: %s", err.Error())
			if silent == false {
				if logToGlog == true {
					glog.Errorln(errorStr)
				} else {
					fmt.Printf("%s", aurora.BrightRed(errorStr))
				}
			}

			return nil
		}
	}

	return v
}

// Check if the configuration file is ok
func IsConfigurationOk(v *viper.Viper, silent bool, logToGlog bool) bool {

	return true
}
