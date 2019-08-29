package xdb

import (
	"github.com/dgraph-io/badger"
	"github.com/golang/glog"
)

type StorageEngine struct {
	db *badger.DB
}

func (se *StorageEngine) Open(dbPath string) bool {
	var err error
	se.db, err = badger.Open(badger.DefaultOptions(dbPath))

	if err != nil {
		glog.Errorln("Failed to open database '", dbPath, "'")
		return false
	}

	return true
}

func (se *StorageEngine) Close() {
	defer se.db.Close()
}

func (se *StorageEngine) Set(key []byte, value []byte) bool {

	err := se.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})

	if err == nil {
		return true
	} else {
		glog.Errorln("Failed to set key:", key, " value:", value)
		return false
	}
}

func (se *StorageEngine) Get(key []byte) ([]byte, bool) {
	var valCopy []byte

	err := se.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return err
	})

	if err == nil {
		return valCopy, true
	} else {
		return valCopy, false
	}
}

func (se *StorageEngine) Delete(key []byte) bool {
	err := se.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})

	if err == nil {
		return true
	} else {
		glog.Errorln("Failed to delete key:", key)
		return false
	}
}
