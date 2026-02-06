package utils

import "crypto/sha256"

func GenerateHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return string(hash.Sum(nil))
}
