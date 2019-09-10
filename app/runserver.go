package app

import (
	"github.com/spf13/viper"
)

func RunServer(v *viper.Viper) {
	hc := PrepareServer(v)
	if hc != nil {
		hc.Setup()
		hc.Start()
	}
}
