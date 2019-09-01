package core

type AuthInterface interface {
	// Get Unique ID from authorization data and also validity of auth data
	GetUID(authorizationData []byte, bearer string) (string, bool)
}
