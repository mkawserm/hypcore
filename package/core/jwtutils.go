package core

import (
	"github.com/pascaldekloe/jwt"
	"strings"
	"time"
)

func tokenFromHeader(auth []byte, bearer string) []byte {
	if string(auth) == "" {
		return []byte("")
	}

	trimmed_bearer := strings.Trim(bearer, " ")

	if !strings.HasPrefix(string(auth), trimmed_bearer+" ") {
		return []byte("")
	}
	return []byte(strings.TrimPrefix(string(auth), trimmed_bearer+" "))
}

func Verify(token []byte, pem []byte, bearer string, isHMACSecret bool) ([]byte, bool) {
	keys := jwt.KeyRegister{}
	if !isHMACSecret {
		_, err := keys.LoadPEM(pem, []byte(""))
		if err != nil {
			return nil, false
		}
	} else {
		keys.Secrets = append(keys.Secrets, pem)
	}

	claims, err2 := keys.Check(tokenFromHeader(token, bearer))

	if err2 == nil {
		if claims.Valid(time.Now()) {
			return claims.Raw, true
		}
	}

	return nil, false
}
