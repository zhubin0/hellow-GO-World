package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

//Calculate block hash.
func CalcBlockHash(block []byte) (string, error) {
	h1 := sha1.New()
	_, err := h1.Write(block)
	if err != nil {
		return "", err
	}
	h256 := sha256.New()
	_, err = h256.Write(block)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h1.Sum(nil)) + hex.EncodeToString(h256.Sum(nil)), nil
}
