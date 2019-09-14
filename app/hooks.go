package app

import (
	"github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/xcore"
	"github.com/mkawserm/hypcore/package/xdb"
	"github.com/spf13/viper"
)

var IsConfigurationOkHook = func(v *viper.Viper, silent bool, logToGlog bool) bool {
	return true
}

var AuthVerifyHook = func() core.AuthVerifyInterface {
	return &core.AuthVerify{}
}

var ServeWSHook = func() core.ServeWSInterface {
	return &core.ServeWSGraphQL{}
}

var OnlineUserDataStoreHook = func() core.OnlineUserDataStoreInterface {
	return core.NewOnlineUserMemoryMap()
}

var StorageEngineHook = func() core.StorageInterface {
	return &xdb.StorageEngine{}
}

var HypCoreSetupHook = func(v *viper.Viper, hypCore *xcore.HypCore, silent bool, logToGlog bool) bool {
	return true
}
