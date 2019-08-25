package core

type ServeWSInterface interface {
	ServeWS(context *HContext, connectionId int, message []byte)
}
