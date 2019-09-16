package core

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"net"
	"net/http"
	"sync"
)

type HContext struct {
	Host string // read only
	Port string // read only

	EventQueueSize int //read only
	WaitingTime    int //read only

	EnableTLS bool   //read only
	CertFile  string //read only
	KeyFile   string //read only

	ServerMux *http.ServeMux

	ConnectionEventPool *EventPool

	IsLive bool          // read write
	Lock   *sync.RWMutex // mutex for modifiable params

	AuthVerify AuthVerifyInterface // read only
	ServeWS    ServeWSInterface    // read only

	AuthBearer     string   // read only
	AuthPublicKey  string   // read only
	AuthPrivateKey string   // read only
	AuthSecretKey  string   // read only
	AuthAlgorithm  string   // read only
	AuthIssuer     string   // read only
	AuthAudiences  []string // read only

	AuthTokenDefaultTimeout      int64 // seconds
	AuthTokenSuperGroupTimeout   int64 // seconds
	AuthTokenServiceGroupTimeout int64 // seconds
	AuthTokenNormalGroupTimeout  int64 // seconds

	OnlineUserDataStore OnlineUserDataStoreInterface // read only

	LivePath             []byte // read only
	AuthPath             []byte // read only
	GraphQLPath          []byte // read only
	WebSocketUpgradePath []byte // read only

	MiddlewareList []MiddlewareInterface // read only while running the server
	RouteList      []*Route              // read only while running the server

	KeyValueStore map[string]string
	StorageEngine StorageInterface
	DbPath        string

	//GraphQL related Objects
	//GraphQLQuery *graphql.Object
	//GraphQLMutation *graphql.Object
	GraphQLQueryFields    graphql.Fields
	GraphQLMutationFields graphql.Fields
	GraphQLSchema         *graphql.Schema

	AuthQueryFields    graphql.Fields
	AuthMutationFields graphql.Fields
	AuthSchema         *graphql.Schema

	EnableLivePath      bool
	EnableAuthPath      bool
	EnableGraphQLPath   bool
	EnableWebSocketPath bool
}

func (c *HContext) AddGraphQLQueryField(name string, field *graphql.Field) {
	c.GraphQLQueryFields[name] = field
}

func (c *HContext) AddGraphQLMutationField(name string, field *graphql.Field) {
	c.GraphQLMutationFields[name] = field
}

func (c *HContext) AddAuthQueryField(name string, field *graphql.Field) {
	c.AuthQueryFields[name] = field
}

func (c *HContext) AddAuthMutationField(name string, field *graphql.Field) {
	c.AuthMutationFields[name] = field
}

func (c *HContext) AddConnection(conn net.Conn) error {
	// Already protected by Mutex
	return c.ConnectionEventPool.AddConnection(conn)
}

func (c *HContext) RemoveConnection(conn net.Conn) error {
	// Already protected by Mutex
	return c.ConnectionEventPool.RemoveConnection(conn)
}

func (c *HContext) GetConnection(cid int) (net.Conn, bool) {
	// Already protected by Mutex
	return c.ConnectionEventPool.GetConnection(cid)
}

func (c *HContext) SetIsLive(live bool) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.IsLive = live
}

func (c *HContext) GetIsLive() bool {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	return c.IsLive
}

func (c *HContext) HasAuthVerify() bool {
	return c.AuthVerify != nil
}

func (c *HContext) HasWSServer() bool {
	return c.ServeWS != nil
}

//func (c *HContext) SetWebSocketUID(webSocketAuth interfaces.AuthVerifyInterface)  {
//	c.AuthVerifyInterface = webSocketAuth
//}

func (c *HContext) AddUser(uid string, group string, sid int) {
	if c.OnlineUserDataStore != nil {
		c.OnlineUserDataStore.AddUser(uid, group, sid)
	}
}

func (c *HContext) RemoveUser(sid int) {
	if c.OnlineUserDataStore != nil {
		c.OnlineUserDataStore.RemoveUser(sid)
	}
}

func (c *HContext) GetIdList(uid string) []int {
	if c.OnlineUserDataStore != nil {
		return c.OnlineUserDataStore.GetIdList(uid)
	} else {
		return make([]int, 0)
	}
}

func (c *HContext) GetUIDFromSID(sid int) string {
	if c.OnlineUserDataStore != nil {
		return c.OnlineUserDataStore.GetUIDFromSID(sid)
	} else {
		return ""
	}
}

func (c *HContext) GetGroupFromSID(sid int) string {
	if c.OnlineUserDataStore != nil {
		return c.OnlineUserDataStore.GetGroupFromSID(sid)
	} else {
		return ""
	}
}

func (c *HContext) WriteMessage(connectionId int, message []byte) {
	conn, ok := c.GetConnection(connectionId)
	if ok {
		err := wsutil.WriteServerText(conn, message)
		if err != nil {
			glog.Errorf("WriteMessage: Failed to write message for ID [%d]\n", connectionId)
			glog.Errorf("WriteMessage: Error message [%s]\n", err)
		}
	} else {
		glog.Errorf("WriteMessage Failed to find user with ID [%d]\n", connectionId)
	}
}

func (c *HContext) WriteLowLevelMessage(connectionId int, opCode ws.OpCode, message []byte) {
	conn, ok := c.GetConnection(connectionId)
	if ok {
		err := wsutil.WriteServerMessage(conn, opCode, message)
		if err != nil {
			glog.Errorf("WriteMessage: Failed to write message for ID [%d]\n", connectionId)
			glog.Errorf("WriteMessage: Error message [%s]\n", err)
		}
	} else {
		glog.Errorf("WriteMessage: Failed to find user with ID [%d]\n", connectionId)
	}
}

func (c *HContext) AddMiddleware(mi MiddlewareInterface) {
	c.MiddlewareList = append(c.MiddlewareList, mi)
}

func (c *HContext) AddRoute(pattern string, httpHandlerObject ServeHTTPInterface) {
	c.RouteList = append(c.RouteList, &Route{
		Pattern:           pattern,
		HttpHandlerObject: httpHandlerObject})
}

func (c *HContext) GetValue(key string) string {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	value, ok := c.KeyValueStore[key]
	if ok {
		return value
	} else {
		return ""
	}
}

func (c *HContext) SetValue(key string, value string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.KeyValueStore[key] = value
}

func (c *HContext) RemoveValue(key string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	delete(c.KeyValueStore, key)
}

func (c *HContext) ClearKeyValueStore() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.KeyValueStore = make(map[string]string)
}

func (c *HContext) SetKeyValueStore(dataMap map[string]string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.KeyValueStore = dataMap
}

func (c *HContext) TotalActiveWebSocketConnections() int {
	return c.ConnectionEventPool.TotalActiveWebSocketConnections()
}

func (c *HContext) SaveToStorage(key []byte, value []byte) bool {
	if c.StorageEngine == nil {
		return false
	} else {
		return c.StorageEngine.Set(key, value)
	}
}

func (c *HContext) GetFromStorage(key []byte) ([]byte, bool) {
	if c.StorageEngine == nil {
		return []byte(""), false
	} else {
		return c.StorageEngine.Get(key)
	}
}

func (c *HContext) DeleteFromStorage(key []byte) bool {
	if c.StorageEngine == nil {
		return false
	} else {
		return c.StorageEngine.Delete(key)
	}
}

func (c *HContext) IsExistsInStorage(key []byte) bool {
	if c.StorageEngine == nil {
		return false
	} else {
		return c.StorageEngine.IsExists(key)
	}
}

func (c *HContext) IsStorageEngineReady() bool {
	return c.StorageEngine != nil
}

func (c *HContext) AddObject(obj interface{}) bool {
	if c.StorageEngine == nil {
		return false
	}

	key := GetPk(obj)
	if key == "" {
		return false
	} else {
		data, err := json.Marshal(obj)
		if err != nil {
			return false
		} else {
			return c.StorageEngine.Set([]byte(key), data)
		}
	}
}

func (c *HContext) GetObject(obj interface{}) bool {
	if c.StorageEngine == nil {
		return false
	}

	key := GetPk(obj)
	if key == "" {
		return false
	} else {
		data, ok := c.StorageEngine.Get([]byte(key))
		if ok {
			err := json.Unmarshal(data, obj)
			if err == nil {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

func (c *HContext) DeleteObject(obj interface{}) bool {
	if c.StorageEngine == nil {
		return false
	}

	key := GetPk(obj)
	if key == "" {
		return false
	} else {
		return c.StorageEngine.Delete([]byte(key))
	}
}

func (c *HContext) IsObjectExists(obj interface{}) bool {
	if c.StorageEngine == nil {
		return false
	}

	key := GetPk(obj)
	return c.StorageEngine.IsExists([]byte(key))
}
