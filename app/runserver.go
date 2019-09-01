package app

import (
	"github.com/mkawserm/hypcore/package/xcore"
	"github.com/spf13/viper"
)

func RunServer(v *viper.Viper) {
	hypCore := xcore.NewHypCore(
		&xcore.HypCoreConfig{
			Host: v.GetString("server.host"),
			Port: v.GetString("server.port"),

			EventQueueSize: v.GetInt("server.eventQueueSize"),
			WaitingTime:    v.GetInt("server.waitingTime"),

			EnableTLS: v.GetBool("server.tls"),
			CertFile:  v.GetString("server.certFile"),
			KeyFile:   v.GetString("server.keyFile"),

			EnableLivePath:      v.GetBool("server.enableLivePath"),
			EnableGraphQLPath:   v.GetBool("server.enableGraphQLPath"),
			EnableWebSocketPath: v.GetBool("server.enableWebSocketPath"),

			DbPath: v.GetString("db.path"),

			Auth:                AuthHook(),
			ServeWS:             ServeWSHook(),
			OnlineUserDataStore: OnlineUserDataStoreHook(),
			StorageEngine:       StorageEngineHook(),
		},
	)

	hypCore.ReconfigurePath(
		[]byte(v.GetString("server.webSocketPath")),
		[]byte(v.GetString("server.graphQLPath")),
		[]byte(v.GetString("server.livePath")),
	)

	if HypCoreSetupHook(v, hypCore, false, true) {
		hypCore.Setup()
		hypCore.Start()
	}
}
