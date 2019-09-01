package core

type StorageInterface interface {
	Open(dbPath string) bool
	Close()
	Set(key []byte, value []byte) bool
	Get(key []byte) ([]byte, bool)
	Delete(key []byte) bool

	IsExists(key []byte) bool
}
