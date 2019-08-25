package core

import (
	"sort"
	"sync"
)

// Thread safe One One Map
type OMap struct {
	aFirst  []interface{}
	aSecond []interface{}

	aLock sync.RWMutex
}

type iterator struct {
	currentIndex int
	omap         *OMap
}

func (i *iterator) Next() bool {
	if i.currentIndex >= i.omap.Length() || i.omap.Length() == 0 {
		return false
	}

	if i.currentIndex == -1 {
		i.currentIndex = 0
		return true
	} else if i.currentIndex < i.omap.Length() {
		i.currentIndex = i.currentIndex + 1
		return true
	} else {
		return false
	}
}

func (i *iterator) Get() (interface{}, interface{}) {
	return i.omap.Get(i.currentIndex)
}

func (o *OMap) New(first []interface{}, second []interface{}) *OMap {
	o.aFirst = first
	o.aSecond = second
	o.aLock = sync.RWMutex{}
	return o
}

func Iterator(o *OMap) *iterator {
	return &iterator{currentIndex: -1, omap: o}
}

func (o *OMap) IsCompatible() bool {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	return len(o.aFirst) == len(o.aSecond)
}

func (o *OMap) GetByFirst(value interface{}) interface{} {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	index := sort.Search(len(o.aFirst), func(i int) bool {
		return o.aFirst[i] == value
	})

	if index >= len(o.aFirst) || o.aFirst[index] != value || index >= len(o.aSecond) {
		return nil
	} else {
		return o.aSecond[index]
	}
}

func (o *OMap) GetBySecond(value interface{}) interface{} {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	index := sort.Search(len(o.aSecond), func(i int) bool {
		return o.aSecond[i] == value
	})

	if index >= len(o.aFirst) || o.aFirst[index] != value || index >= len(o.aSecond) {
		return nil
	} else {
		return o.aFirst[index]
	}
}

func (o *OMap) Append(key, value interface{}) {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	o.aFirst = append(o.aFirst, key)
	o.aSecond = append(o.aSecond, value)
}

func (o *OMap) DeleteByFirst(value interface{}) bool {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	index := sort.Search(len(o.aFirst), func(i int) bool {
		return o.aFirst[i] == value
	})

	if index >= len(o.aFirst) || o.aFirst[index] != value || index >= len(o.aSecond) {
		return false
	} else {
		o.aFirst[index] = o.aFirst[len(o.aFirst)-1]
		o.aFirst = o.aFirst[0 : len(o.aFirst)-1]

		o.aSecond[index] = o.aSecond[len(o.aSecond)-1]
		o.aSecond = o.aSecond[0 : len(o.aSecond)-1]

		return true
	}
}

func (o *OMap) DeleteBySecond(value interface{}) bool {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	index := sort.Search(len(o.aSecond), func(i int) bool {
		return o.aSecond[i] == value
	})

	if index >= len(o.aFirst) || o.aFirst[index] != value || index >= len(o.aSecond) {
		return false
	} else {
		o.aFirst[index] = o.aFirst[len(o.aFirst)-1]
		o.aFirst = o.aFirst[0 : len(o.aFirst)-1]

		o.aSecond[index] = o.aSecond[len(o.aSecond)-1]
		o.aSecond = o.aSecond[0 : len(o.aSecond)-1]

		return true
	}
}

func (o *OMap) DeleteByIndex(index int) bool {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	if index < 0 || index >= len(o.aFirst) || index >= len(o.aSecond) {
		return false
	} else {
		o.aFirst[index] = o.aFirst[len(o.aFirst)-1]
		o.aFirst = o.aFirst[0 : len(o.aFirst)-1]

		o.aSecond[index] = o.aSecond[len(o.aSecond)-1]
		o.aSecond = o.aSecond[0 : len(o.aSecond)-1]

		return true
	}
}

func (o *OMap) GetFirstByIndex(index int) interface{} {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	if index < 0 || index >= len(o.aFirst) {
		return nil
	} else {
		return o.aFirst[index]
	}
}

func (o *OMap) GetSecondByIndex(index int) interface{} {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	if index < 0 || index >= len(o.aFirst) {
		return nil
	} else {
		return o.aSecond[index]
	}
}

func (o *OMap) Length() int {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	if len(o.aFirst) == len(o.aSecond) {
		return len(o.aFirst)
	} else {
		return 0
	}
}

func (o *OMap) Get(index int) (interface{}, interface{}) {
	o.aLock.Lock()
	defer o.aLock.Unlock()

	if len(o.aFirst) == 0 || len(o.aSecond) == 0 || index < 0 {
		return nil, nil
	} else if len(o.aFirst) == len(o.aSecond) {
		return o.aFirst[index], o.aSecond[index]
	} else {
		return nil, nil
	}
}
