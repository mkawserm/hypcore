package xdb

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/golang/glog"
)

type badgerLogger struct {
}

func (bl *badgerLogger) Errorf(f string, v ...interface{}) {
	glog.Errorf("ERROR: "+f+"\n", v...)
}

func (bl *badgerLogger) Warningf(f string, v ...interface{}) {
	glog.Warningf("WARNING: "+f+"\n", v...)
}

func (bl *badgerLogger) Infof(f string, v ...interface{}) {
	glog.Infof("INFO: "+f+"\n", v...)
}

func (bl *badgerLogger) Debugf(f string, v ...interface{}) {
	glog.Infof("DEBUG: "+f+"\n", v...)
}

func StorageEngineOptions(path string) badger.Options {
	return badger.Options{
		Dir:                 path,
		ValueDir:            path,
		LevelOneSize:        256 << 20,
		LevelSizeMultiplier: 10,
		TableLoadingMode:    options.MemoryMap,
		ValueLogLoadingMode: options.MemoryMap,
		// table.MemoryMap to mmap() the tables.
		// table.Nothing to not preload the tables.
		MaxLevels:               7,
		MaxTableSize:            64 << 20,
		NumCompactors:           2, // Compactions can be expensive. Only run 2.
		NumLevelZeroTables:      5,
		NumLevelZeroTablesStall: 10,
		NumMemtables:            5,
		SyncWrites:              true,
		NumVersionsToKeep:       1,
		CompactL0OnClose:        true,
		// Nothing to read/write value log using standard File I/O
		// MemoryMap to mmap() the value log files
		// (2^30 - 1)*2 when mmapping < 2^31 - 1, max int32.
		// -1 so 2*ValueLogFileSize won't overflow on 32-bit systems.
		ValueLogFileSize: 1<<30 - 1,

		ValueLogMaxEntries: 1000000,
		ValueThreshold:     32,
		Truncate:           false,
		Logger:             &badgerLogger{},
		LogRotatesToFlush:  2,
	}
}

type StorageEngine struct {
	db        *badger.DB
	gcChannel chan bool
}

func (se *StorageEngine) Open(dbPath string) bool {
	se.gcChannel = make(chan bool, 1000)

	var err error
	se.db, err = badger.Open(StorageEngineOptions(dbPath))

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
		if err == nil {
			err = item.Value(func(val []byte) error {
				valCopy = append([]byte{}, val...)
				return nil
			})
		}
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
