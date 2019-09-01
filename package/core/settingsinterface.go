package core

type SettingsInterface interface {
	GetSettings(key string) []byte
	SetSettings(key string, value []byte) bool
}
