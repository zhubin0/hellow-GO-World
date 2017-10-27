package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

// Sha512Hash provides Sha512 hash functionality.
func Sha512Hash(data string, salt string) []byte {
	h := sha512.New()
	if salt != "" {
		h.Write([]byte(salt))
	}
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Sha256Hash provides Sha256 hash functionality.
func Sha256Hash(data string, salt string) []byte {
	h := sha256.New()
	if salt != "" {
		h.Write([]byte(salt))
	}
	h.Write([]byte(data))
	return h.Sum(nil)
}

//HexSha512Hash provides Sha512 hash functionality,and return hash in lower-case hex.
func HexSha512Hash(data string, salt string) string {
	return hex.EncodeToString(Sha512Hash(data, salt))
}

//HexSha256Hash provides Sha256 hash functionality,and return hash in lower-case hex.
func HexSha256Hash(data string, salt string) string {
	return hex.EncodeToString(Sha256Hash(data, salt))
}

// Sha3Hash provides Sha3 hash functionality.
func Sha3Hash(data string) []byte {
	buf := []byte(data)
	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	h := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(h, buf)
	return h
}

// Sha1Hash provides Sha1 hash functionality.
func Sha1Hash(data string) []byte {
	h := sha1.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Md5Hash provides MD5 hash functionality.
func Md5Hash(data string) []byte {
	b := []byte(data)
	hmd5 := md5.New()
	hmd5.Write(b)
	return hmd5.Sum(nil)
}

// Md5Sha1Hash implements hybrid hash function which consists of the
// concatenation of an MD5 and SHA1 hash. Use this if you want to reduce possible
// MD5 collision possibility. Return value is in Hex.
func Md5Sha1Hash(data string) string {
	b := []byte(data)
	md5sha1 := make([]byte, md5.Size+sha1.Size)
	hsha1 := sha1.New()
	hsha1.Write(b)
	hmd5 := md5.New()
	hmd5.Write(b)
	copy(md5sha1, hmd5.Sum(nil))
	copy(md5sha1[md5.Size:], hsha1.Sum(nil))
	return hex.EncodeToString(md5sha1)
}

// Sha3HashBase64Safe wraps the Sha3 hash into a safe encoded base64 string.
func Sha3HashBase64Safe(data string) (base64safe string) {
	return base64.RawURLEncoding.EncodeToString(Sha3Hash(data))
}

// Generate md5 string
func GenMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Generate md5 string
func GenBytesToMd5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// MD5Encode encodes the given data using md5, returning the hexadecimal results.
func MD5Encode(data string) string {
	b := md5.Sum([]byte(data))
	//return fmt.Sprintf("%x", b)
	return hex.EncodeToString(b[:])
}
