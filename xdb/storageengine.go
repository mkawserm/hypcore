package xdb

import (
	"github.com/dgraph-io/badger"
	"github.com/golang/glog"
)

type StorageEngine struct {
	DbPath string

	db *badger.DB
}

func (se *StorageEngine) Open() bool {
	var err error
	se.db, err = badger.Open(badger.DefaultOptions(se.DbPath))

	if err != nil {
		glog.Errorln("Failed to open database '", se.DbPath, "'")
		return false
	}

	return true
}

func (se *StorageEngine) Close() {
	defer se.db.Close()
}
