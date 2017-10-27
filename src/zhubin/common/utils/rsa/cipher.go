package rsa

import (
	"bytes"
	"crypto"
	"fmt"
	"math/big"

	"errors"

	"code.google.com/p/log4go"
)

var MAX_DECRYPT_BLOCK = 128

type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
	PubKeyDecrypt(cipherText []byte) ([]byte, error)
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

func NewCipher(key Key, padding Padding, cipherMode CipherMode, signMode SignMode) Cipher {
	return &cipher{key: key, padding: padding, cipherMode: cipherMode, sign: signMode}
}

type cipher struct {
	key        Key
	cipherMode CipherMode
	sign       SignMode
	padding    Padding
}

// Encrypt with public key.
func (cipher *cipher) Encrypt(plainText []byte) ([]byte, error) {
	groups := cipher.padding.Padding(plainText)
	buffer := bytes.Buffer{}
	for _, plainTextBlock := range groups {
		cipherText, err := cipher.cipherMode.Encrypt(plainTextBlock, cipher.key.PublicKey())
		if err != nil {
			log4go.Error(err)
			return nil, err
		}
		buffer.Write(cipherText)
	}
	return buffer.Bytes(), nil
}

// Decrypt with private key.
func (cipher *cipher) Decrypt(cipherText []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, errors.New("encrypted data can not be nil")
	}
	/*
		BUG记录：传入的cipherText为空数组时，则会导致解密失败，因此对数据分组的算法要仔细检查。
	*/
	groups := grouping(cipherText, cipher.key.Modulus())
	buffer := bytes.Buffer{}
	for _, cipherTextBlock := range groups {
		plainText, err := cipher.cipherMode.Decrypt(cipherTextBlock, cipher.key.PrivateKey())
		if err != nil {
			log4go.Error(err)
			return nil, err
		}
		buffer.Write(plainText)
	}
	return buffer.Bytes(), nil
}

// Decrypt with public key.
func (cipher *cipher) PubKeyDecrypt(data []byte) ([]byte, error) {
	enT := data
	inputLen := len(enT)
	offSet := 0
	i := 0
	cache := make([]byte, 0)
	out := make([]byte, 0)
	var err error

	for { //TODO error handle.
		if inputLen-offSet > MAX_DECRYPT_BLOCK {
			cache, err = cipher.pubKeyDecrypt(enT[offSet : MAX_DECRYPT_BLOCK+offSet])
		} else {
			if offSet == 0 {
				cache, err = cipher.pubKeyDecrypt(enT)
			} else {
				cache, err = cipher.pubKeyDecrypt(enT[offSet:inputLen])
			}
		}
		if err != nil {
			return nil, err
		}
		for i, _ := range cache {
			out = append(out, cache[i])
		}
		cache = cache[:0]
		i++
		offSet = i * MAX_DECRYPT_BLOCK
		if inputLen-offSet <= 0 {
			break
		}
	}
	return out, nil
}

func (cipher *cipher) pubKeyDecrypt(data []byte) ([]byte, error) {
	pub := cipher.key.PublicKey()
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, fmt.Errorf("Invalid data length.")
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, fmt.Errorf("Data too large.")
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, fmt.Errorf("Data is broken.")
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, fmt.Errorf("Key pair dismatch.")
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

func (cipher *cipher) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	return cipher.sign.Sign(src, hash, cipher.key.PrivateKey())
}

func (cipher *cipher) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	return cipher.sign.Verify(src, sign, hash, cipher.key.PublicKey())
}

//copy from 'crypto/rsa'
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}
