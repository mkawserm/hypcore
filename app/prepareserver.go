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

		AuthTokenDefaultTimeout:      v.GetUint32("auth.tokenDefaultTimeout"),
		AuthTokenNormalGroupTimeout:  v.GetUint32("auth.tokenNormalGroupTimeout"),
		AuthTokenServiceGroupTimeout: v.GetUint32("auth.tokenServiceGroupTimeout"),
		AuthTokenSuperGroupTimeout:   v.GetUint32("auth.tokenSuperGroupTimeout"),
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
