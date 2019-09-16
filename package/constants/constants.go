package constants

import "net/textproto"

var (
	//headerHost          = "Host"
	headerUpgrade       = "Upgrade"
	headerConnection    = "Connection"
	headerAuthorization = "Authorization"
	//headerSecVersion    = "Sec-WebSocket-Version"
	headerSecProtocol = "Sec-WebSocket-Protocol"
	//headerSecExtensions = "Sec-WebSocket-Extensions"
	//headerSecKey        = "Sec-WebSocket-Key"
	//headerSecAccept     = "Sec-WebSocket-Accept"

	//headerHostCanonical          = textproto.CanonicalMIMEHeaderKey(headerHost)
	HeaderUpgradeCanonical       = textproto.CanonicalMIMEHeaderKey(headerUpgrade)
	HeaderConnectionCanonical    = textproto.CanonicalMIMEHeaderKey(headerConnection)
	HeaderAuthorizationCanonical = textproto.CanonicalMIMEHeaderKey(headerAuthorization)
	//headerSecVersionCanonical    = textproto.CanonicalMIMEHeaderKey(headerSecVersion)
	HeaderSecProtocolCanonical = textproto.CanonicalMIMEHeaderKey(headerSecProtocol)
	//headerSecExtensionsCanonical = textproto.CanonicalMIMEHeaderKey(headerSecExtensions)
	//headerSecKeyCanonical        = textproto.CanonicalMIMEHeaderKey(headerSecKey)
	//headerSecAcceptCanonical     = textproto.CanonicalMIMEHeaderKey(headerSecAccept)
)

const (
	SuperGroup   = "super"
	ServiceGroup = "service"
	NormalGroup  = "normal"
)

const WebSocketConnection = 1
const WebSocketSecureConnection = 2

const HTTPConnection = 3
const HTTPSConnection = 4

const (
	WebSocketOpContinuation byte = 0x0
	WebSocketOpText         byte = 0x1
	WebSocketOpBinary       byte = 0x2
	WebSocketOpClose        byte = 0x8
	WebSocketOpPing         byte = 0x9
	WebSocketOpPong         byte = 0xa
)
