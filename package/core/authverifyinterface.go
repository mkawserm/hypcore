package core

type AuthVerifyInterface interface {
	// Get Unique ID from authorization data and also validity of auth data
	Verify(token []byte, authPublicKey []byte, authSecretKey []byte, authBearer string) (map[string]interface{}, bool)
}
