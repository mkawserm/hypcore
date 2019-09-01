package core

import (
	"sort"
	"sync"
)

type OnlineUserMemoryMap struct {
	UserMap           map[string][]int
	SocketIDToUserMap map[int]string
	UserMapLock       *sync.RWMutex
}

func NewOnlineUserMemoryMap() *OnlineUserMemoryMap {
	return &OnlineUserMemoryMap{
		UserMap:           make(map[string][]int),
		SocketIDToUserMap: make(map[int]string),
		UserMapLock:       &sync.RWMutex{},
	}
}

func (o *OnlineUserMemoryMap) AddUser(uid string, sid int) {
	o.UserMapLock.Lock()
	defer o.UserMapLock.Unlock()

	o.SocketIDToUserMap[sid] = uid

	if _, ok := o.UserMap[uid]; ok {
		o.UserMap[uid] = append(o.UserMap[uid], sid)
	} else {
		o.UserMap[uid] = make([]int, 1)
		o.UserMap[uid] = append(o.UserMap[uid], sid)
	}
}

func (o *OnlineUserMemoryMap) RemoveUser(sid int) {
	o.UserMapLock.Lock()
	defer o.UserMapLock.Unlock()

	uid, ok1 := o.SocketIDToUserMap[sid]

	if ok1 {
		delete(o.SocketIDToUserMap, sid)

		val, ok := o.UserMap[uid]

		if ok {
			index := sort.Search(len(val), func(i int) bool {
				return val[i] == sid
			})

			if index >= len(val) || val[index] != sid {
				return
			} else {
				val[index] = val[len(val)-1]
				val = val[0 : len(val)-1]
			}

			if len(val) == 0 {
				delete(o.UserMap, uid)
			}
		}
	}
}

func (o *OnlineUserMemoryMap) GetIdList(uid string) []int {
	o.UserMapLock.RLock()
	defer o.UserMapLock.RUnlock()

	val, ok := o.UserMap[uid]

	if ok {
		tmp := make([]int, len(val))
		copy(tmp, val)
		return tmp
	} else {
		return make([]int, 0)
	}
}

func (o *OnlineUserMemoryMap) GetUIDFromSID(sid int) string {
	o.UserMapLock.RLock()
	defer o.UserMapLock.RUnlock()

	val, ok := o.SocketIDToUserMap[sid]
	if ok {
		return val
	} else {
		return ""
	}
}
