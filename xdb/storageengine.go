package xdb

import (
	"github.com/dgraph-io/badger"
	"github.com/golang/glog"
)

type StorageEngine struct {
	db        *badger.DB
	gcChannel chan bool
}

func (se *StorageEngine) Open(dbPath string) bool {
	se.gcChannel = make(chan bool, 1000)

	var err error
	se.db, err = badger.Open(badger.DefaultOptions(dbPath))

	if err != nil {
		glog.Errorln("Failed to open database '", dbPath, "'")
		return false
	}

	go func() {
		select {
		case msg := <-se.gcChannel:
			if msg {
				se.runGc()
			}
		}
	}()

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
		se.gcChannel <- true
		return true
	} else {
		glog.Errorln("Failed to delete key:", key)
		return false
	}
}

func (se *StorageEngine) IsExists(key []byte) bool {
	_, b := se.Get(key)
	return b
}

func (se *StorageEngine) runGc() {
	_ = se.db.RunValueLogGC(0.7)
	//ticker := time.NewTicker(5 * time.Minute)
	//defer ticker.Stop()
	//for range ticker.C {
	//again:
	//	err := se.db.RunValueLogGC(0.7)
	//	if err == nil {
	//		goto again
	//	}
	//}
}
