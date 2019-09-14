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

func IsHMAC(token []byte) bool {
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

func Verify(token []byte, keyData []byte, bearer string) ([]byte, bool) {
	keys := jwt.KeyRegister{}
	if !IsHMAC(token) {
		_, err := keys.LoadPEM(keyData, []byte(""))
		if err != nil {
			glog.Errorf("LoadPEM failed: %s\n", err.Error())
			return nil, false
		}
	} else {
		keys.Secrets = append(keys.Secrets, keyData)
	}

	claims, err2 := keys.Check(tokenFromHeader(token, bearer))

	if err2 == nil {
		if claims.Valid(time.Now()) {
			return claims.Raw, true
		} else {
			glog.Errorln("JWT Time expired\n")
		}
	} else {
		glog.Errorf("JWT Signature check failed: %s\n", err2.Error())
	}

	return nil, false
}
