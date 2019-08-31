package app

import (
	"github.com/mkawserm/hypcore/core"
	"github.com/mkawserm/hypcore/xcore"
	"github.com/spf13/viper"
)

var IsConfigurationOkHook = func(v *viper.Viper, silent bool, logToGlog bool) bool {
	return true
}

var AuthHook = func() core.AuthInterface {
	return nil
}

var ServeWSHook = func() core.ServeWSInterface {
	return nil
}

var OnlineUserDataStoreHook = func() core.OnlineUserDataStoreInterface {
	return nil
}

var StorageEngineHook = func() core.StorageInterface {
	return nil
}

var HypCoreSetupHook = func(hypCore *xcore.HypCore, silent bool, logToGlog bool) bool {
	return true
}
