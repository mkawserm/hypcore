package core

type ServeWSInterface interface {
	ServeWS(ctx *HContext, connectionId int, message []byte, opCode byte)
}
