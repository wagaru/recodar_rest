package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandToken(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
