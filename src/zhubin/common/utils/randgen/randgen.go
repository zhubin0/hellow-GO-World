package randgen

import (
	"crypto/rand"
	"datamesh.com/common/utils/base62"
	"datamesh.com/common/utils/flake"
	"datamesh.com/common/utils/uuid"
)

//generate uuid
func GenUUID() string {
	return uuid.NewV5(uuid.NewV4(), "mesh-expert").String()
}

// =================================== random string =================================================

// base62
var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func GenRandString(n int) string {
	return genRandString(n, stdChars)
}

// =================================== random Mongo ID =================================================

var flk *flake.Flake

func init() {
	f, err := flake.New()
	if err != nil {
		panic(err)
	}
	flk = f
}

// generate a mongodb ID using flake, in base62 encoding.
// a 64bit flake plus a random 88bit noise
// the total char count for the generated id is 26 or 27
// NOTE todo should test on K8S
func GenMongoId() string {
	id := flk.NextId()
	b := make([]byte, 11)
	safeRandom(b)
	return base62.EncodeToString(append(id.Bytes(), b...))
}

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

// =================================== unique string =================================================

// we make this unique string based on flake and random numbers, in base62 encoding.
// there is no way to make sure everything is unique when we use local random string generator with
// arbitrary machine (k8s, vm, etc).
// so we have to make this string long enough so we can safely say it is unique (collision possibility is extremely low)
// padding now is byte number instead of base64 character count: the prefix length is 26 or 27 chars (long enough), you may
// consider add few padding
func GenUniqueString(padding int) string {
	id := flk.NextId()
	b := make([]byte, 11+padding)
	safeRandom(b)
	return base62.EncodeToString(append(id.Bytes(), b...))
}

// =================================== random number =================================================

var numChars = []byte("0123456789")

func GenRandNumString(n int) string {
	return genRandString(n, numChars)
}

func genRandString(length int, chars []byte) string {
	clen := len(chars)
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

func GenRandBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}
