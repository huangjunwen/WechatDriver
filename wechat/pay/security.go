package pay

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"sort"
)

// New MD5 hash instance.
func newMD5() func() hash.Hash {

	return md5.New

}

// New HMAC-SHA256 hash instance.
func newHMACSHA256(key string) func() hash.Hash {

	return func() hash.Hash {

		return hmac.New(sha256.New, []byte(key))

	}

}

// Sign a collection of data (dict) and return hex digest in lower case.
// See: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
func signDict(dict map[string]string, new_hash func() hash.Hash, hash_key string) string {

	// Sort keys
	keys := make(sort.StringSlice, 0, len(dict))

	for k, _ := range dict {

		keys = append(keys, k)

	}

	sort.Sort(keys)

	// Hash query string like string "param1=value1&param2=value2&...&key=key"
	h := new_hash()

	for _, k := range keys {

		v := dict[k]

		// Skip empty value
		if v == "" {

			continue

		}

		h.Write([]byte(k))

		h.Write([]byte("="))

		h.Write([]byte(v))

		h.Write([]byte("&"))

	}

	h.Write([]byte("key="))

	h.Write([]byte(hash_key))

	// return the hex digest
	return fmt.Sprintf("%x", h.Sum(nil))

}
