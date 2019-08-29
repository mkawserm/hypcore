package core

type OnlineUserDataStoreInterface interface {
	AddUser(uid string, sid int)
	RemoveUser(sid int)
	GetIdList(uid string) []int
	GetUIDFromSID(sid int) string
}
