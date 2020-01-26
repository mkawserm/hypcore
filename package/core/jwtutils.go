package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/pascaldekloe/jwt"
	"strings"
	"time"
)

// Header is a critical subset of the registered “JOSE Header Parameter Names”.
type header struct {
	Alg  string   // algorithm
	Kid  string   // key identifier
	Crit []string // extensions which must be understood and processed
}

var b64encoding = base64.RawURLEncoding

func IsHMACAlg(token []byte) bool {
	firstDot := bytes.IndexByte(token, '.')
	lastDot := bytes.LastIndexByte(token, '.')
	if lastDot <= firstDot {
		return false
	}
	buf := make([]byte, b64encoding.DecodedLen(len(token)))
	n, err := b64encoding.Decode(buf, token[:firstDot])
	if err != nil {
		return false
	}

	h := new(header)
	if err := json.Unmarshal(buf[:n], h); err != nil {
		return false
	}

	switch h.Alg {
	case jwt.HS256:
		return true
	case jwt.HS384:
		return true
	case jwt.HS512:
		return true
	default:
		return false
	}
}

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

func VerifyJWT(token []byte,
	authPublicKey []byte,
	authSecretKey []byte,
	authBearer string,
	plainToken bool,
	logToGlog bool) (map[string]interface{}, bool) {

	keys := jwt.KeyRegister{}
	newToken := []byte("")

	if plainToken {
		newToken = token
	} else {
		newToken = tokenFromHeader(token, authBearer)
	}

	if !IsHMACAlg(newToken) {
		_, err := keys.LoadPEM([]byte(authPublicKey), []byte(""))
		if err != nil {
			if logToGlog {
				glog.Errorf("LoadPEM failed: %s\n", err.Error())
			}

			return nil, false
		}
	} else {
		keys.Secrets = append(keys.Secrets, []byte(authSecretKey))
	}

	var claims *jwt.Claims
	var err2 error

	claims, err2 = keys.Check(newToken)

	if err2 == nil {
		// glog.Infoln(string(claims.Raw))
		if claims.Valid(time.Now()) {
			data := make(map[string]interface{})
			err3 := json.Unmarshal(claims.Raw, &data)
			if err3 == nil {
				return data, true
			} else {
				if logToGlog {
					glog.Errorf("JSON Unmarshal error: %s\n", err3.Error())
				}

				return nil, false
			}
		} else {
			if logToGlog {
				glog.Errorln("JWT Time expired")
			}
		}
	} else {
		if logToGlog {
			glog.Errorf("JWT Signature check failed: %s\n", err2.Error())
		}
	}

	return nil, false
}
