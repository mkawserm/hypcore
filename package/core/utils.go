package core

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/gobwas/httphead"
	"reflect"
	"unsafe"
)

const (
	toLower = 'a' - 'A' // for use with OR.
	//toUpper  = ^byte(toLower) // for use with AND.
	toLower8 = uint64(toLower) |
		uint64(toLower)<<8 |
		uint64(toLower)<<16 |
		uint64(toLower)<<24 |
		uint64(toLower)<<32 |
		uint64(toLower)<<40 |
		uint64(toLower)<<48 |
		uint64(toLower)<<56
)

func StrToBytes(str string) (bts []byte) {
	s := (*reflect.StringHeader)(unsafe.Pointer(&str))
	b := (*reflect.SliceHeader)(unsafe.Pointer(&bts))
	b.Data = s.Data
	b.Len = s.Len
	b.Cap = s.Len
	return
}

func StrHasToken(header, token string) (has bool) {
	return BtsHasToken(StrToBytes(header), StrToBytes(token))
}

func BtsHasToken(header, token []byte) (has bool) {
	httphead.ScanTokens(header, func(v []byte) bool {
		has = BtsEqualFold(v, token)
		return !has
	})
	return
}

// StrEqualFold checks s to be case insensitive equal to p.
// Note that p must be only ascii letters. That is, every byte in p belongs to
// range ['a','z'] or ['A','Z'].
func StrEqualFold(s, p string) bool {
	return BtsEqualFold(StrToBytes(s), StrToBytes(p))
}

// btsEqualFold checks s to be case insensitive equal to p.
// Note that p must be only ascii letters. That is, every byte in p belongs to
// range ['a','z'] or ['A','Z'].
func BtsEqualFold(s, p []byte) bool {
	if len(s) != len(p) {
		return false
	}
	n := len(s)
	// Prepare manual conversion on bytes that not lay in uint64.
	m := n % 8
	for i := 0; i < m; i++ {
		if s[i]|toLower != p[i]|toLower {
			return false
		}
	}
	// Iterate over uint64 parts of s.
	n = (n - m) >> 3
	if n == 0 {
		// There are no more bytes to compare.
		return true
	}

	for i := 0; i < n; i++ {
		x := m + (i << 3)
		av := *(*uint64)(unsafe.Pointer(&s[x]))
		bv := *(*uint64)(unsafe.Pointer(&p[x]))
		if av|toLower8 != bv|toLower8 {
			return false
		}
	}

	return true
}

func IsObjectStructType(obj interface{}) bool {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return t.Elem().Kind() == reflect.Struct
	} else {
		return t.Kind() == reflect.Struct
	}
}

func GetObjectTypeName(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return t.Elem().PkgPath() + "::" + t.Elem().Name()
	} else {
		return t.Elem().PkgPath() + "::" + t.Name()
	}
}

func GetPk(obj interface{}) string {
	if IsObjectStructType(obj) {
		typeName := GetObjectTypeName(obj)
		elementsField := reflect.ValueOf(obj).Elem()
		pk := elementsField.FieldByName("Pk")
		if pk.IsValid() && pk.Kind() == reflect.String {
			keyName := string("<" + typeName + "::Pk::" + pk.String() + ">")
			return keyName
		}
	}

	return ""
}

func LoadPrivateKey(data []byte) (interface{}, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("Invalid pem block")
	}

	switch block.Type {
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		return key, err

	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		return key, err

	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		return key, err
	}

	return nil, errors.New("Failed to load private key")
}

func ParseGraphQLData(gqlData []byte) (string, map[string]interface{}, error) {
	query := ""
	variables := make(map[string]interface{})

	var queryMap map[string]interface{}

	err := json.Unmarshal(gqlData, &queryMap)
	if err == nil {
		if q, ok := queryMap["query"]; ok {
			query = q.(string)
		} else {
			return "", variables, errors.New("no query key found")
		}

		if v, ok := queryMap["variables"]; ok {
			variables = v.(map[string]interface{})
		}

	} else {
		return "", nil, err
	}

	if query == "" {
		return "", nil, errors.New("empty graphql query string")
	} else {
		return query, variables, nil
	}
}
