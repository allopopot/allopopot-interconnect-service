package utility

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecret(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.RawStdEncoding.EncodeToString(bytes)
}
