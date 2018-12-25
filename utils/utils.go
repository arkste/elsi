package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// CreateHashFromString create a md5 hash of text
func CreateHashFromString(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
