package rsa

import (
	"encoding/base64"
)

type rsa_helper struct {
}

func NewHelper() *rsa_helper {
	return &rsa_helper{}
}

//decrypt by rsa public key.
func (rh *rsa_helper) DecryptByPublicKey(publicKey string, cipherText string) ([]byte, error) {
	return rh.decrypt(true, publicKey, cipherText)
}

//decrypt by rsa private key.
func (rh *rsa_helper) DecryptByPrivateKey(privateKey, cipherText string) ([]byte, error) {
	return rh.decrypt(false, privateKey, cipherText)
}

//encrypt by rsa public key to Base64 text.
func (rh *rsa_helper) EncryptByPublicKey(publicKey string, plainText string) (string, error) {
	return rh.encrypt(true, publicKey, plainText)
}

//encrypt by rsa private key to Base64 text.
func (rh *rsa_helper) EncryptByPrivateKey(privateKey, plainText string) (string, error) {
	return rh.encrypt(false, privateKey, plainText)
}

//encrypt plain text.
//if isPub is true,key represent public key, key represents private key.
func (rh *rsa_helper) encrypt(isPub bool, key string, plainText string) (string, error) {
	cipher, err := rh.newCipher(isPub, key)
	if err != nil {
		return "", err
	}
	enT, err := cipher.Encrypt([]byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enT), nil
}

//decrypt cipherText by key.
//if isPub is true,key represent public key, key represents private key.
func (rh *rsa_helper) decrypt(isPub bool, key string, cipherText string) ([]byte, error) {
	cipher, err := rh.newCipher(isPub, key)
	if err != nil {
		return nil, err
	}
	enT, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	if isPub {
		return cipher.PubKeyDecrypt(enT)
	}
	return cipher.Decrypt(enT)

}

//create cipher by key.
func (rh *rsa_helper) newCipher(isPub bool, key string) (cipher Cipher, err error) {
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	var _key Key
	if isPub {
		_key, err = ParsePKCS8PubKey(k)
	} else {
		_key, err = ParsePKCS8PriKey(k)
	}
	if err != nil {
		return nil, err
	}
	cipher, err = NewRSA(_key)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}
