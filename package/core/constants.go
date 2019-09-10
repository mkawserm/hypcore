package core

import "net/textproto"

var (
	//headerHost          = "Host"
	headerUpgrade       = "Upgrade"
	headerConnection    = "Connection"
	headerAuthorization = "Authorization"
	//headerSecVersion    = "Sec-WebSocket-Version"
	//headerSecProtocol   = "Sec-WebSocket-Protocol"
	//headerSecExtensions = "Sec-WebSocket-Extensions"
	//headerSecKey        = "Sec-WebSocket-Key"
	//headerSecAccept     = "Sec-WebSocket-Accept"

	//headerHostCanonical          = textproto.CanonicalMIMEHeaderKey(headerHost)
	HeaderUpgradeCanonical       = textproto.CanonicalMIMEHeaderKey(headerUpgrade)
	HeaderConnectionCanonical    = textproto.CanonicalMIMEHeaderKey(headerConnection)
	HeaderAuthorizationCanonical = textproto.CanonicalMIMEHeaderKey(headerAuthorization)
	//headerSecVersionCanonical    = textproto.CanonicalMIMEHeaderKey(headerSecVersion)
	//headerSecProtocolCanonical   = textproto.CanonicalMIMEHeaderKey(headerSecProtocol)
	//headerSecExtensionsCanonical = textproto.CanonicalMIMEHeaderKey(headerSecExtensions)
	//headerSecKeyCanonical        = textproto.CanonicalMIMEHeaderKey(headerSecKey)
	//headerSecAcceptCanonical     = textproto.CanonicalMIMEHeaderKey(headerSecAccept)
)

const (
	SuperGroup   = "super"
	ServiceGroup = "service"
	NormalGroup  = "normal"
)
