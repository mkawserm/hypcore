package core

type OnlineUserDataStoreInterface interface {
	AddUser(uid string, group string, sid int)
	RemoveUser(sid int)
	GetIdList(uid string) []int
	GetUIDFromSID(sid int) string
	GetGroupFromSID(sid int) string

	GetUIDList() []string
	GetTotalUID() int
}
