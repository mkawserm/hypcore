package app

import (
	"github.com/mkawserm/hypcore/z"
	"github.com/spf13/viper"
	"strings"
)

var ConfigFilePathFirst = "/home/" + strings.ToLower(z.ExeName())
var ConfigFileNameWithoutExt = strings.ToLower(z.ExeName())

//func ConfigFilePathSecond() string {
//	return "/etc/" + strings.ToLower(z.ExeName())
//}

func ConfigFileName() string {
	return strings.ToLower(z.ExeName()) + ".toml"
}

func NewConfigFile(configFilePath string, configFileName string) *viper.Viper {
	if configFilePath != "" {
		ConfigFilePathFirst = configFilePath
		if ConfigFilePathFirst[len(ConfigFilePathFirst)-1:] == "/" {
			ConfigFilePathFirst = ConfigFilePathFirst[0 : len(ConfigFilePathFirst)-1]
		}
	}

	if configFileName != "" {
		ConfigFileNameWithoutExt = configFileName
	}

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName(ConfigFileNameWithoutExt)
	v.AddConfigPath(ConfigFilePathFirst)

	//v.AddConfigPath("/etc/"+strings.ToLower(z.ExeName()))

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")

	v.SetDefault("server.tls", false)
	v.SetDefault("server.certFile", ConfigFilePathFirst+"/cert.pem")
	v.SetDefault("server.keyFile", ConfigFilePathFirst+"/key.pem")

	v.SetDefault("server.enableLivePath", true)
	v.SetDefault("server.enableGraphQLPath", true)
	v.SetDefault("server.enableWebSocketPath", true)

	v.SetDefault("server.livePath", "/api/live")
	v.SetDefault("server.graphQLPath", "/graphql")
	v.SetDefault("server.webSocketPath", "/ws")

	v.SetDefault("server.eventQueueSize", 100)
	v.SetDefault("server.waitingTime", 100)

	v.SetDefault("db.path", ConfigFilePathFirst+"/db")

	return v
}
