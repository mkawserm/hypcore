package core

type StorageInterface interface {
	Open(dbPath string) bool
	Close()

	Set(key []byte, value []byte) bool
	Get(key []byte) ([]byte, bool)
	Delete(key []byte) bool

	IsExists(key []byte) bool

	AddObject(obj interface{}) bool
	GetObject(obj interface{}) bool
	DeleteObject(obj interface{}) bool
	IsObjectExists(obj interface{}) bool
}
