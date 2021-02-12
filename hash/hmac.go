package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// HMAC wrapper around crypto/hmac
type HMAC struct {
	hmac hash.Hash
}

// NewHMAC returns a new HMAC object
func NewHMAC(key string) HMAC {
	hm := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: hm,
	}
}

// Hash returns the  hash string
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
