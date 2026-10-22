package tokenhash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func Compare(hash, value string) bool {
	return Hash(value) == hash
}
