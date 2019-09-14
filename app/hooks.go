package app

import (
	"github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/xcore"
	"github.com/spf13/viper"
)

var IsConfigurationOkHook = func(v *viper.Viper, silent bool, logToGlog bool) bool {
	return true
}

var AuthVerifyHook = func() core.AuthVerifyInterface {
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

var HypCoreSetupHook = func(v *viper.Viper, hypCore *xcore.HypCore, silent bool, logToGlog bool) bool {
	return true
}
