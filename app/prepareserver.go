package app

import (
	"github.com/mkawserm/hypcore/package/xcore"
	"github.com/spf13/viper"
)

func PrepareServer(v *viper.Viper) *xcore.HypCore {
	hcc := &xcore.HypCoreConfig{
		Host: v.GetString("server.host"),
		Port: v.GetString("server.port"),

		EventQueueSize: v.GetInt("server.eventQueueSize"),
		WaitingTime:    v.GetInt("server.waitingTime"),

		EnableTLS: v.GetBool("server.tls"),
		CertFile:  v.GetString("server.certFile"),
		KeyFile:   v.GetString("server.keyFile"),

		EnableAuthVerify: v.GetBool("server.enableAuthVerify"),

		EnableLivePath:      v.GetBool("server.enableLivePath"),
		EnableAuthPath:      v.GetBool("server.enableAuthPath"),
		EnableGraphQLPath:   v.GetBool("server.enableGraphQLPath"),
		EnableWebSocketPath: v.GetBool("server.enableWebSocketPath"),

		DbPath: v.GetString("db.path"),

		AuthVerify:          AuthVerifyHook(),
		ServeWS:             ServeWSHook(),
		OnlineUserDataStore: OnlineUserDataStoreHook(),
		StorageEngine:       StorageEngineHook(),

		AuthBearer:     v.GetString("auth.bearer"),
		AuthAlgorithm:  v.GetString("auth.algorithm"),
		AuthPublicKey:  v.GetString("auth.publicKey"),
		AuthPrivateKey: v.GetString("auth.privateKey"),
		AuthSecretKey:  v.GetString("auth.secretKey"),

		AuthTokenDefaultTimeout:      v.GetInt64("auth.tokenDefaultTimeout"),
		AuthTokenNormalGroupTimeout:  v.GetInt64("auth.tokenNormalGroupTimeout"),
		AuthTokenServiceGroupTimeout: v.GetInt64("auth.tokenServiceGroupTimeout"),
		AuthTokenSuperGroupTimeout:   v.GetInt64("auth.tokenSuperGroupTimeout"),

		CORSAllowAllOrigins:  v.GetBool("cors.AllowAllOrigins"),
		CORSAllowCredentials: v.GetBool("cors.AllowCredentials"),
		CORSMaxAge:           v.GetDuration("cors.MaxAge"),

		CORSAllowOrigins:  v.GetStringSlice("cors.AllowOrigins"),
		CORSAllowMethods:  v.GetStringSlice("cors.AllowMethods"),
		CORSAllowHeaders:  v.GetStringSlice("cors.AllowHeaders"),
		CORSExposeHeaders: v.GetStringSlice("cors.ExposeHeaders"),
	}

	hypCore := xcore.NewHypCore(hcc)

	hypCore.ReconfigurePath(
		[]byte(v.GetString("server.webSocketPath")),
		[]byte(v.GetString("server.graphQLPath")),
		[]byte(v.GetString("server.livePath")),
		[]byte(v.GetString("server.authPath")),
	)

	if HypCoreSetupHook(v, hypCore, false, true) {
		return hypCore
	}
	return nil
}
