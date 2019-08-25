package hypcore

import (
	"flag"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/mkawserm/hypcore/core"
	"github.com/mkawserm/hypcore/views"
	"net/http"
	"sync"
	"syscall"
)

type HypCore struct {
	context *core.HContext
	ready   bool
}

type HypCoreConfig struct {
	Host           string
	Port           string
	EventQueueSize int
	WaitingTime    int

	Auth                core.AuthInterface
	ServeWS             core.ServeWSInterface
	OnlineUserDataStore core.OnlineUserDataStoreInterface

	EnableTLS bool
	CertFile  string
	KeyFile   string
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
	}

	hContext := &core.HContext{
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

		WebSocketUpgradePath: []byte("/ws"),
		GraphQLPath:          []byte("/graphql"),
		LivePath:             []byte("/api/live"),

		KeyValueStore:         make(map[string]string),
		GraphQLQueryFields:    make(graphql.Fields),
		GraphQLMutationFields: make(graphql.Fields),
	}

	if hContext.ServeWS == nil {
		hContext.ServeWS = &core.ServeWSGraphQL{}
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
	h.AddGraphQLQueryField("isLive", &graphql.Field{
		Type:        graphql.Boolean,
		Resolve:     func(p graphql.ResolveParams) (interface{}, error) { return h.context.GetIsLive(), nil },
		Description: "Check if the service is live",
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

	h.context.ConnectionEventPool, err = core.MakeCustomEventPool(
		h.context.EventQueueSize,
		h.context.WaitingTime)

	if err != nil {
		panic(err)
	}

	h.context.ServerMux = http.NewServeMux()

	h.context.ServerMux.Handle(string(h.context.LivePath),
		&views.LiveView{
			Context: h.context,
		})

	h.context.ServerMux.Handle(string(h.context.WebSocketUpgradePath),
		&views.WebSocketUpgradeView{
			Context: h.context,
		})

	h.context.ServerMux.Handle(string(h.context.GraphQLPath),
		&views.GraphQLView{
			Context: h.context,
		})

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
					h.context.RemoveUser(core.WebsocketFileDescriptor(conn))
				}

				if err := h.context.RemoveConnection(conn); err != nil {
					glog.Infof("Failed to remove %v", err)
				}

				_ = conn.Close()

			} else {
				// Actual Message
				//log.Printf("msg: %s", string(msg))

				// Call Message Server to process the message
				if h.context.ServeWS != nil {
					h.context.ServeWS.ServeWS(h.context, core.WebsocketFileDescriptor(conn), msg)
				}
			}
		}
	}
}

// Start HypCore server
func (h *HypCore) Start() {
	if h.ready == false {
		panic("please call Setup() first")
	}
	if h.context == nil {
		panic("context has not properly configured.")
	}

	if h.context.ConnectionEventPool == nil {
		panic("Connection event pool has not properly configured.")
	}

	if h.context.ServerMux == nil {
		panic("Server Mux has not properly configured.")
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

func (h *HypCore) AddMiddleware(mi core.MiddlewareInterface) {
	h.context.AddMiddleware(mi)
}

func (h *HypCore) AddRoute(pattern string, httpHandlerObject core.ServeHTTPInterface) {
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
