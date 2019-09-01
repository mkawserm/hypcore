package core

import (
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
