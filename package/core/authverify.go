package core

type AuthVerify struct {
}

func (av *AuthVerify) Verify(token []byte, authPublicKey []byte, authSecretKey []byte, authBearer string) (map[string]interface{}, bool) {
	return VerifyJWT(token, authPublicKey, authSecretKey, authBearer, false, true)
}
