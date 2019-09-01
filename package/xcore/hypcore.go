package xcore

import (
	"flag"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/views"
	xdb2 "github.com/mkawserm/hypcore/package/xdb"
	"net/http"
	"runtime"
	"sync"
	"syscall"
)

type HypCore struct {
	context *core2.HContext
	ready   bool
}

type HypCoreConfig struct {
	Host           string
	Port           string
	EventQueueSize int
	WaitingTime    int

	EnableTLS bool
	CertFile  string
	KeyFile   string

	EnableAuth bool

	EnableLivePath      bool
	EnableGraphQLPath   bool
	EnableWebSocketPath bool

	DbPath string

	AuthBearer     string
	AuthPublicKey  string
	AuthPrivateKey string
	AuthAlgorithm  string

	Auth                core2.AuthInterface
	ServeWS             core2.ServeWSInterface
	OnlineUserDataStore core2.OnlineUserDataStoreInterface
	StorageEngine       core2.StorageInterface
}

func init() {
	//flag.Usage =
	flag.Parse()
}

//Create new HypCore server
func NewHypCore(hc *HypCoreConfig) *HypCore {
	if hc.Auth != nil && hc.OnlineUserDataStore == nil {
		glog.Fatal("Auth found but no OnlineUserDataStore found. Please configure OnlineUserDataStore.")
		return nil
	} else if hc.Auth == nil && hc.OnlineUserDataStore != nil {
		glog.Fatal("OnlineUserDataStore found but no Auth found. Please configure Auth.")
		return nil
	} else if hc.Auth != nil && hc.OnlineUserDataStore != nil {
		if hc.AuthBearer == "" {
			glog.Fatal("AuthBearer is required but not provided")
		}

		if hc.AuthAlgorithm == "" {
			glog.Fatal("AuthAlgorithm is required but not provided")
		}

		if hc.AuthPublicKey == "" {
			glog.Fatal("AuthPublicKey is required but not provided")
		}

		if hc.AuthPrivateKey == "" {
			glog.Fatal("AuthPrivateKey is required but not provided")
		}
	}

	hContext := &core2.HContext{
		Host:           hc.Host,
		Port:           hc.Port,
		EventQueueSize: hc.EventQueueSize,
		WaitingTime:    hc.WaitingTime,

		EnableTLS: hc.EnableTLS,
		CertFile:  hc.CertFile,
		KeyFile:   hc.KeyFile,

		ServerMux:           nil,
		ConnectionEventPool: nil,

		Lock: &sync.RWMutex{},

		IsLive: true,

		Auth:                hc.Auth,
		ServeWS:             hc.ServeWS,
		OnlineUserDataStore: hc.OnlineUserDataStore,

		StorageEngine: hc.StorageEngine,
		DbPath:        hc.DbPath,

		WebSocketUpgradePath: []byte("/ws"),
		GraphQLPath:          []byte("/graphql"),
		LivePath:             []byte("/api/live"),

		KeyValueStore:         make(map[string]string),
		GraphQLQueryFields:    make(graphql.Fields),
		GraphQLMutationFields: make(graphql.Fields),

		EnableLivePath:      hc.EnableLivePath,
		EnableGraphQLPath:   hc.EnableGraphQLPath,
		EnableWebSocketPath: hc.EnableWebSocketPath,

		AuthBearer:     hc.AuthBearer,
		AuthPublicKey:  hc.AuthPublicKey,
		AuthAlgorithm:  hc.AuthAlgorithm,
		AuthPrivateKey: hc.AuthPrivateKey,
	}

	if hContext.ServeWS == nil {
		hContext.ServeWS = &core2.ServeWSGraphQL{}
	}

	h := &HypCore{
		context: hContext,
		ready:   false,
	}
	return h
}

func (h *HypCore) ReconfigurePath(webSocketUpgradePath []byte, graphQLPath []byte, livePath []byte) {
	h.context.WebSocketUpgradePath = webSocketUpgradePath
	h.context.GraphQLPath = graphQLPath
	h.context.LivePath = livePath
}

func (h *HypCore) Setup() {
	if h.context.StorageEngine == nil {
		h.context.StorageEngine = &xdb2.StorageEngine{}
	}

	if h.context.ServeWS == nil {
		h.context.ServeWS = &core2.ServeWSGraphQL{}
	}

	if h.context.StorageEngine != nil {
		h.context.StorageEngine.Open(h.context.DbPath)
		runtime.SetFinalizer(h, func(h *HypCore) {
			h.context.StorageEngine.Close()
		})
	}

	h.AddGraphQLQueryField("isLive", &graphql.Field{
		Type:        graphql.Boolean,
		Resolve:     func(p graphql.ResolveParams) (interface{}, error) { return h.context.GetIsLive(), nil },
		Description: "Check if the service is live",
	})

	h.AddGraphQLQueryField("totalActiveWebSocketConnections", &graphql.Field{
		Type: graphql.Int,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return h.context.TotalActiveWebSocketConnections(), nil
		},
		Description: "Query total active websocket connections",
	})

	h.AddGraphQLMutationField("updateLive", &graphql.Field{
		Type: graphql.Boolean,
		Args: graphql.FieldConfigArgument{
			"live": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			live := params.Args["live"].(bool)
			h.context.SetIsLive(live)
			return h.context.GetIsLive(), nil
		},
		Description: "Update service availability",
	})

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:        "Query",
			Fields:      h.context.GraphQLQueryFields,
			Description: "GraphQL Query",
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:        "Mutation",
			Fields:      h.context.GraphQLMutationFields,
			Description: "GraphQL Mutation",
		}),
	})

	h.context.GraphQLSchema = schema

	if h.context.EventQueueSize == 0 {
		h.context.EventQueueSize = 100
	}

	if h.context.WaitingTime == 0 {
		h.context.WaitingTime = 100
	}

	var err error

	h.context.ConnectionEventPool, err = core2.MakeCustomEventPool(
		h.context.EventQueueSize,
		h.context.WaitingTime)

	if err != nil {
		panic(err)
	}

	h.context.ServerMux = http.NewServeMux()

	if h.context.EnableLivePath == true {
		h.context.ServerMux.Handle(string(h.context.LivePath),
			&views.LiveView{
				Context: h.context,
			})
	}

	if h.context.EnableWebSocketPath == true {
		h.context.ServerMux.Handle(string(h.context.WebSocketUpgradePath),
			&views.WebSocketUpgradeView{
				Context: h.context,
			})
	}

	if h.context.EnableGraphQLPath == true {
		h.context.ServerMux.Handle(string(h.context.GraphQLPath),
			&views.GraphQLView{
				Context: h.context,
			})
	}

	h.context.ServerMux.Handle(string([]byte("/")),
		&views.DynamicView{
			Context: h.context,
		})

	h.ready = true
}

func (h *HypCore) runMainEventLoop() {
	for {
		connections, err := h.context.ConnectionEventPool.Wait()
		if err != nil {
			glog.Warningf("Failed to wait on eventPool %v", err)
			continue
		}

		for _, conn := range connections {
			if conn == nil {
				break
			}

			if msg, _, err := wsutil.ReadClientData(conn); err != nil {
				if h.context.HasAuth() {
					h.context.RemoveUser(core2.WebsocketFileDescriptor(conn))
				}

				if err := h.context.RemoveConnection(conn); err != nil {
					glog.Infof("Failed to remove %v", err)
				}

				_ = conn.Close()

			} else {
				// Actual Message
				// glog.Infoln("msg: %s", string(msg))

				// Call Message Server to process the message
				if h.context.ServeWS != nil {
					h.context.ServeWS.ServeWS(h.context, core2.WebsocketFileDescriptor(conn), msg)
				}
			}
		}
	}
}

// Start HypCore server
func (h *HypCore) Start() {
	if h.ready == false {
		glog.Fatalln("please call Setup() first")
	}
	if h.context == nil {
		glog.Fatalln("context has not properly configured.")
	}

	if h.context.ConnectionEventPool == nil {
		glog.Fatalln("Connection event pool has not properly configured.")
	}

	if h.context.ServerMux == nil {
		glog.Fatalln("Server Mux has not properly configured.")
	}

	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	glog.Infoln("HypCore Server started.")

	go h.runMainEventLoop()

	if h.context.EnableTLS {
		glog.Infoln("Server is listening at: https://" + h.context.Host + ":" + h.context.Port)
		err := http.ListenAndServeTLS(h.context.Host+":"+h.context.Port,
			h.context.CertFile,
			h.context.KeyFile, h.context.ServerMux)
		if err != nil {
			glog.Infoln("http TLS server error: ", err)
		}

	} else {
		glog.Infoln("Server is listening at: http://" + h.context.Host + ":" + h.context.Port)
		err := http.ListenAndServe(h.context.Host+":"+h.context.Port, h.context.ServerMux)
		if err != nil {
			glog.Infoln("http server error: ", err)
		}
	}

	glog.Infoln("HypCore Server finished.")
}

// Set if the HypCore server is live or not
func (h *HypCore) SetIsLive(live bool) {
	h.context.SetIsLive(live)
}

func (h *HypCore) HasAuth() bool {
	return h.context.HasAuth()
}

func (h *HypCore) AddMiddleware(mi core2.MiddlewareInterface) {
	h.context.AddMiddleware(mi)
}

func (h *HypCore) AddRoute(pattern string, httpHandlerObject core2.ServeHTTPInterface) {
	h.context.AddRoute(pattern, httpHandlerObject)
}

func (h *HypCore) AddGraphQLQueryField(name string, field *graphql.Field) {
	h.context.AddGraphQLQueryField(name, field)
}

func (h *HypCore) AddGraphQLMutationField(name string, field *graphql.Field) {
	h.context.AddGraphQLMutationField(name, field)
}

// Get value using the key from the key value store of context
func (h *HypCore) GetValue(key string) string {
	return h.context.GetValue(key)
}

// Set value in the key value store of context
func (h *HypCore) SetValue(key string, value string) {
	h.context.SetValue(key, value)
}

// Remove a value from the key value store of context
func (h *HypCore) RemoveValue(key string) {
	h.context.RemoveValue(key)
}

// Clear key value store of context
func (h *HypCore) ClearKeyValueStore() {
	h.context.ClearKeyValueStore()
}

func (h *HypCore) SetKeyValueStore(dataMap map[string]string) {
	h.context.SetKeyValueStore(dataMap)
}
