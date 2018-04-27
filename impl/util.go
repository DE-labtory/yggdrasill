package impl

import "crypto/sha256"

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}
