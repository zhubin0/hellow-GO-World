package cryptos

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
)

// AESEncryptHex works with AESDecryptHexVary as a pair. It encrypts the given
// text using AES, and produces a hexadecimal encoded text.
func AESEncryptHexVary(key []byte, text string) (string, error) {
	b, err := AESEncrypt(key, []byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// AESDecryptHex works with AESEncryptHexVary as a pair. It decrypts the given
// text ciphered by AES, and return the original plain text.
func AESDecryptHexVary(key []byte, text string) (string, error) {
	h, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	b, err := AESDecrypt(key, h)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// AESEncryptHex works with AESDecryptHex as a pair. It encrypts the given
// text using AES, and produces a hexadecimal encoded text.
func AESEncryptHex(key []byte, text string) (string, error) {
	b, err := AesEncrypt(key, []byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// AESDecryptHex works with AESEncryptHex as a pair. It decrypts the given
// text ciphered by AES, and return the original plain text.
func AESDecryptHex(key []byte, text string) (string, error) {
	h, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	b, err := AesDecrypt(key, h)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func AesEncrypt(key, origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(key, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Base64EncodeSafe encodes a byte slice into a URL safe string and gets rid of the padding.
func Base64EncodeSafe(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64DecodeSafe works as the counter-part of Base64EncodeSafe to decode a URL-safe
// base64 string to byte slice.
func Base64DecodeSafe(text string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(text)
}

// AESEncryptBase64Safe works with AESDecryptBase64Safe as a pair. It encrypts the given
// text using AES, and produces a base64-encoded text which is safe to use in URL or Cookie.
func AESEncryptBase64Safe(key []byte, text string) (string, error) {
	b, err := AESEncrypt(key, []byte(text))
	if err != nil {
		return "", err
	}
	return Base64EncodeSafe(b), nil
}

// AESDecryptBase64Safe works with AESEncryptBase64Safe as a pair. It decrypts the given
// text ciphered by AES, and return the original plain text.
func AESDecryptBase64Safe(key []byte, text string) (string, error) {
	h, err := Base64DecodeSafe(text)
	if err != nil {
		return "", err
	}
	b, err := AESDecrypt(key, h)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// AESEncrypt performs the basic AES encryption.
// NOTE it uses random padding which produces different encrypted text each time, be careful
func AESEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	return ciphertext, nil
}

// AESDecrypt performs the basic AES decryption.
func AESDecrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return text, nil
}

func RSAEncryptBase64Safe(key []byte, data string) (string, error) {
	b, err := RSAEncrypt(key, []byte(data))
	if err != nil {
		return "", err
	}
	return Base64EncodeSafe(b), nil
}

func RSADecryptBase64Safe(key []byte, data string) (string, error) {
	b, err := Base64DecodeSafe(data)
	if err != nil {
		return "", err
	}
	p, err := RSADecrypt(key, []byte(b))
	if err != nil {
		return "", err
	}
	return string(p), nil
}

func RSAEncrypt(key, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RSADecrypt(key, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}