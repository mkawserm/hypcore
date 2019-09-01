package app

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"os"
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
	var e error

	if v.IsSet("db.path") {
		dbPath := v.GetString("db.path")
		e = EnsureDir(dbPath)
		if e != nil {
			if silent == false {
				if logToGlog {
					glog.Errorln(e.Error())
				} else {
					fmt.Println(aurora.BrightRed(e.Error()))
				}
			}
			return false
		}
	} else {
		SpitError("db.path key is missing in the configuration", silent, logToGlog)
		return false
	}

	if v.IsSet("server.tls") && v.GetBool("server.tls") {

		if v.IsSet("server.certFile") {
			_, e1 := os.Stat(v.GetString("server.certFile"))
			if CheckAndSpit(e1, silent, logToGlog) {
				return false
			}
		} else {
			SpitError("server.certFile key is missing in the configuration", silent, logToGlog)
			return false
		}

		if v.IsSet("server.keyFile") {
			_, e1 := os.Stat(v.GetString("server.keyFile"))
			if CheckAndSpit(e1, silent, logToGlog) {
				return false
			}
		} else {
			SpitError("server.keyFile key is missing in the configuration", silent, logToGlog)
			return false
		}

	}

	if v.IsSet("server.enableAuth") && v.GetBool("server.enableAuth") {

		if v.IsSet("auth.bearer") {
			if v.GetString("auth.bearer") == "" {
				SpitError("auth.bearer can not be empty", silent, logToGlog)
				return false
			}
		} else {
			SpitError("auth.bearer key is missing in the configuration", silent, logToGlog)
			return false
		}

		if v.IsSet("auth.algorithm") {
			if v.GetString("auth.algorithm") == "" {
				SpitError("auth.algorithm can not be empty", silent, logToGlog)
				return false
			}
		} else {
			SpitError("auth.algorithm key is missing in the configuration", silent, logToGlog)
			return false
		}

		if v.IsSet("auth.publicKey") {
			if v.GetString("auth.publicKey") == "" {
				SpitError("auth.publicKey can not be empty", silent, logToGlog)
				return false
			}
		} else {
			SpitError("auth.publicKey key is missing in the configuration", silent, logToGlog)
			return false
		}

		if v.IsSet("auth.privateKey") {
			if v.GetString("auth.privateKey") == "" {
				SpitError("auth.privateKey can not be empty", silent, logToGlog)
				return false
			}
		} else {
			SpitError("auth.privateKey key is missing in the configuration", silent, logToGlog)
			return false
		}

	}

	return IsConfigurationOkHook(v, silent, logToGlog)
}

func SpitError(message string, silent bool, logToGlog bool) {
	if silent == false {
		if logToGlog {
			glog.Errorln(message)
		} else {
			fmt.Println(aurora.BrightRed(message))
		}
	}
}

func CheckAndSpit(e error, silent bool, logToGlog bool) bool {
	if e != nil {
		if silent == false {
			if logToGlog {
				glog.Errorln(e.Error())
			} else {
				fmt.Println(aurora.BrightRed(e.Error()))
			}
		}
		// exit
		return true
	}
	// don't exit
	return false
}
